package wallet

import (
	"encoding/json"
	"errors"
	"math/big"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ethAddressRe = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

var ErrInvalidAddress = errors.New("invalid ethereum address")

const tokenCacheTTL = time.Hour

type Labeler interface {
	LabelFor(address string) string
}

type Service struct {
	client *Client
	repo   *Repository
	labels Labeler
}

func NewService(client *Client, repo *Repository, labels Labeler) *Service {
	return &Service{client: client, repo: repo, labels: labels}
}

func (s *Service) Lookup(userID, address string) (*WalletSnapshot, error) {
	if !ethAddressRe.MatchString(address) {
		return nil, ErrInvalidAddress
	}

	weiHex, raw, err := s.client.GetNativeBalance(address)
	if err != nil {
		return nil, err
	}
	wei := parseBalance(weiHex)

	snap := &WalletSnapshot{
		UserID:      userID,
		Address:     address,
		BalanceWei:  wei.String(),
		BalanceEth:  formatUnits(wei, 18),
		RawResponse: raw,
	}
	if err := s.repo.Save(snap); err != nil {
		return nil, err
	}
	return snap, nil
}

func (s *Service) History(userID string) ([]WalletSnapshot, error) {
	return s.repo.ListByUser(userID)
}

func (s *Service) Overview(address string) (*Overview, error) {
	if !ethAddressRe.MatchString(address) {
		return nil, ErrInvalidAddress
	}

	transfers, err := s.client.GetAllTransfers(address)
	if err != nil {
		return nil, err
	}

	var incoming, outgoing int
	first, last := "", ""
	counterparties := map[string]int{}

	for _, t := range transfers {
		ts := t.Metadata.BlockTimestamp
		if ts != "" {
			if first == "" || ts < first {
				first = ts
			}
			if ts > last {
				last = ts
			}
		}

		if strings.EqualFold(t.To, address) {
			incoming++
			if t.From != "" && !strings.EqualFold(t.From, address) {
				counterparties[strings.ToLower(t.From)]++
			}
		} else {
			outgoing++
			if t.To != "" && !strings.EqualFold(t.To, address) {
				counterparties[strings.ToLower(t.To)]++
			}
		}
	}

	ageDays := 0
	if first != "" {
		if ft, perr := time.Parse(time.RFC3339, first); perr == nil {
			ageDays = int(time.Since(ft).Hours() / 24)
		}
	}

	top := topCounterparties(counterparties, 5)
	for i := range top {
		top[i].Label = s.labels.LabelFor(top[i].Address)
	}

	return &Overview{
		Address:           address,
		TxCount:           len(transfers),
		IncomingCount:     incoming,
		OutgoingCount:     outgoing,
		FirstTx:           first,
		LastTx:            last,
		AgeDays:           ageDays,
		TopCounterparties: top,
	}, nil
}

const graphMaxNodes = 20

func (s *Service) Graph(address string) (*Graph, error) {
	if !ethAddressRe.MatchString(address) {
		return nil, ErrInvalidAddress
	}

	transfers, err := s.client.GetAllTransfers(address)
	if err != nil {
		return nil, err
	}

	type stat struct{ in, out int }
	stats := map[string]*stat{}

	for _, t := range transfers {
		if strings.EqualFold(t.To, address) {
			if t.From == "" || strings.EqualFold(t.From, address) {
				continue
			}
			k := strings.ToLower(t.From)
			if stats[k] == nil {
				stats[k] = &stat{}
			}
			stats[k].in++
		} else {
			if t.To == "" || strings.EqualFold(t.To, address) {
				continue
			}
			k := strings.ToLower(t.To)
			if stats[k] == nil {
				stats[k] = &stat{}
			}
			stats[k].out++
		}
	}

	type ranked struct {
		addr string
		s    *stat
	}
	list := make([]ranked, 0, len(stats))
	for a, st := range stats {
		list = append(list, ranked{a, st})
	}
	sort.Slice(list, func(i, j int) bool {
		return (list[i].s.in + list[i].s.out) > (list[j].s.in + list[j].s.out)
	})
	if len(list) > graphMaxNodes {
		list = list[:graphMaxNodes]
	}

	root := strings.ToLower(address)
	nodes := []GraphNode{{ID: root, Label: s.labels.LabelFor(root), Type: "root"}}
	edges := make([]GraphEdge, 0, len(list))

	for _, r := range list {
		direction := "both"
		if r.s.out == 0 {
			direction = "in"
		} else if r.s.in == 0 {
			direction = "out"
		}
		nodes = append(nodes, GraphNode{ID: r.addr, Label: s.labels.LabelFor(r.addr), Type: "counterparty"})
		edges = append(edges, GraphEdge{Source: root, Target: r.addr, Count: r.s.in + r.s.out, Direction: direction})
	}

	return &Graph{Nodes: nodes, Edges: edges}, nil
}

func topCounterparties(m map[string]int, n int) []Counterparty {
	list := make([]Counterparty, 0, len(m))
	for addr, count := range m {
		list = append(list, Counterparty{Address: addr, Count: count})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Count > list[j].Count
	})
	if len(list) > n {
		list = list[:n]
	}
	return list
}

func (s *Service) Transactions(address string) ([]Transaction, error) {
	if !ethAddressRe.MatchString(address) {
		return nil, ErrInvalidAddress
	}

	transfers, err := s.client.GetTransfers(address)
	if err != nil {
		return nil, err
	}

	sort.Slice(transfers, func(i, j int) bool {
		return hexToInt64(transfers[i].BlockNum) > hexToInt64(transfers[j].BlockNum)
	})
	if len(transfers) > 25 {
		transfers = transfers[:25]
	}

	txs := make([]Transaction, 0, len(transfers))
	for _, t := range transfers {
		direction := "out"
		if strings.EqualFold(t.To, address) {
			direction = "in"
		}
		txs = append(txs, Transaction{
			Hash:        t.Hash,
			From:        t.From,
			To:          t.To,
			Value:       strconv.FormatFloat(t.Value, 'f', -1, 64),
			Asset:       t.Asset,
			Timestamp:   t.Metadata.BlockTimestamp,
			BlockNumber: parseBalance(t.BlockNum).String(),
			Direction:   direction,
		})
	}
	return txs, nil
}

const (
	minTokenUSD = 0.01
	maxTokenUSD = 1e10
)

func (s *Service) TokenHoldings(address string, refresh bool) (*Portfolio, error) {
	if !ethAddressRe.MatchString(address) {
		return nil, ErrInvalidAddress
	}

	if !refresh {
		if cached, err := s.repo.GetTokenCache(address); err == nil {
			if time.Since(cached.FetchedAt) < tokenCacheTTL {
				var tokens []TokenHolding
				if json.Unmarshal([]byte(cached.TokensJSON), &tokens) == nil {
					return &Portfolio{TotalUSD: cached.TotalUSD, Tokens: tokens}, nil
				}
			}
		}
	}

	tokens, _, err := s.client.GetPortfolio(address)
	if err != nil {
		return nil, err
	}

	holdings := make([]TokenHolding, 0)
	total := big.NewFloat(0)

	for _, t := range tokens {
		if t.TokenAddress == "" {
			continue
		}
		price := usdPrice(t.TokenPrices)
		if price == "" {
			continue
		}
		bal := parseBalance(t.TokenBalance)
		if bal.Sign() <= 0 {
			continue
		}
		balance := formatUnits(bal, t.TokenMetadata.Decimals)
		val := usdValueFloat(balance, price)
		v, _ := val.Float64()
		if v < minTokenUSD || v > maxTokenUSD {
			continue
		}
		total.Add(total, val)
		holdings = append(holdings, TokenHolding{
			Contract: t.TokenAddress,
			Name:     t.TokenMetadata.Name,
			Symbol:   t.TokenMetadata.Symbol,
			Decimals: t.TokenMetadata.Decimals,
			Balance:  balance,
			USDPrice: price,
			USDValue: val.Text('f', 2),
			Logo:     t.TokenMetadata.Logo,
		})
	}

	sort.Slice(holdings, func(i, j int) bool {
		return parseFloat(holdings[i].USDValue) > parseFloat(holdings[j].USDValue)
	})

	portfolio := &Portfolio{TotalUSD: total.Text('f', 2), Tokens: holdings}

	if data, err := json.Marshal(portfolio.Tokens); err == nil {
		_ = s.repo.SaveTokenCache(&TokenCache{
			Address:    address,
			TotalUSD:   portfolio.TotalUSD,
			TokensJSON: string(data),
			FetchedAt:  time.Now(),
		})
	}

	return portfolio, nil
}

func parseBalance(s string) *big.Int {
	n := new(big.Int)
	str, base := s, 10
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		str, base = s[2:], 16
	}
	if _, ok := n.SetString(str, base); !ok {
		return big.NewInt(0)
	}
	return n
}

func hexToInt64(h string) int64 {
	return parseBalance(h).Int64()
}

func formatUnits(amount *big.Int, decimals int) string {
	if decimals <= 0 {
		return amount.String()
	}
	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	return new(big.Float).Quo(new(big.Float).SetInt(amount), divisor).Text('f', decimals)
}

func usdPrice(prices []alchemyTokenPrice) string {
	for _, p := range prices {
		if p.Currency == "usd" {
			return p.Value
		}
	}
	return ""
}

func usdValueFloat(balance, price string) *big.Float {
	b, _, err1 := big.ParseFloat(balance, 10, 200, big.ToNearestEven)
	p, _, err2 := big.ParseFloat(price, 10, 200, big.ToNearestEven)
	if err1 != nil || err2 != nil {
		return big.NewFloat(0)
	}
	return new(big.Float).Mul(b, p)
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
