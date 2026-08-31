package wallet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Client struct {
	http   *http.Client
	apiKey string
}

func NewClient() *Client {
	_ = godotenv.Load()

	return &Client{
		http:   &http.Client{Timeout: 15 * time.Second},
		apiKey: os.Getenv("ALCHEMY_API_KEY"),
	}
}

func (c *Client) GetPortfolio(address string) ([]alchemyToken, string, error) {
	if c.apiKey == "" {
		return nil, "", errors.New("ALCHEMY_API_KEY not set")
	}

	url := fmt.Sprintf("https://api.g.alchemy.com/data/v1/%s/assets/tokens/by-address", c.apiKey)

	var all []alchemyToken
	firstRaw := ""
	pageKey := ""

	for page := 0; page < 25; page++ {
		reqBody := map[string]any{
			"addresses": []map[string]any{
				{"address": address, "networks": []string{"eth-mainnet"}},
			},
			"withMetadata":        true,
			"withPrices":          true,
			"includeNativeTokens": true,
		}
		if pageKey != "" {
			reqBody["pageKey"] = pageKey
		}

		body, err := c.postJSON(url, reqBody)
		if err != nil {
			return nil, "", err
		}
		if firstRaw == "" {
			firstRaw = string(body)
		}

		var pr alchemyPortfolioResponse
		if err := json.Unmarshal(body, &pr); err != nil {
			return nil, "", err
		}
		all = append(all, pr.Data.Tokens...)

		if pr.Data.PageKey == "" {
			break
		}
		pageKey = pr.Data.PageKey
	}

	return all, firstRaw, nil
}

func (c *Client) GetNativeBalance(address string) (string, string, error) {
	if c.apiKey == "" {
		return "", "", errors.New("ALCHEMY_API_KEY not set")
	}

	url := fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s", c.apiKey)

	reqBody := map[string]any{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"params":  []any{address, "latest"},
	}

	body, err := c.postJSON(url, reqBody)
	if err != nil {
		return "", "", err
	}

	var br alchemyBalanceResponse
	if err := json.Unmarshal(body, &br); err != nil {
		return "", "", err
	}
	return br.Result, string(body), nil
}

func (c *Client) GetTransfers(address string) ([]alchemyTransfer, error) {
	return c.transfers(address, "0x19")
}

func (c *Client) GetAllTransfers(address string) ([]alchemyTransfer, error) {
	return c.transfers(address, "0x3e8")
}

func (c *Client) transfers(address, maxCount string) ([]alchemyTransfer, error) {
	if c.apiKey == "" {
		return nil, errors.New("ALCHEMY_API_KEY not set")
	}

	url := fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s", c.apiKey)

	out, err := c.assetTransfers(url, "fromAddress", address, maxCount)
	if err != nil {
		return nil, err
	}
	in, err := c.assetTransfers(url, "toAddress", address, maxCount)
	if err != nil {
		return nil, err
	}
	return append(out, in...), nil
}

func (c *Client) assetTransfers(url, direction, address, maxCount string) ([]alchemyTransfer, error) {
	params := map[string]any{
		"fromBlock":        "0x0",
		"toBlock":          "latest",
		direction:          address,
		"category":         []string{"external"},
		"withMetadata":     true,
		"excludeZeroValue": false,
		"maxCount":         maxCount,
		"order":            "desc",
	}
	reqBody := map[string]any{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "alchemy_getAssetTransfers",
		"params":  []any{params},
	}

	body, err := c.postJSON(url, reqBody)
	if err != nil {
		return nil, err
	}

	var tr alchemyTransfersResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, err
	}
	return tr.Result.Transfers, nil
}

func (c *Client) postJSON(url string, payload any) ([]byte, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alchemy error %d: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}
