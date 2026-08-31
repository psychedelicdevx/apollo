package wallet

import "time"

type alchemyPortfolioResponse struct {
	Data struct {
		Tokens  []alchemyToken `json:"tokens"`
		PageKey string         `json:"pageKey"`
	} `json:"data"`
}

type alchemyToken struct {
	TokenAddress  string              `json:"tokenAddress"`
	TokenBalance  string              `json:"tokenBalance"`
	TokenMetadata alchemyTokenMeta    `json:"tokenMetadata"`
	TokenPrices   []alchemyTokenPrice `json:"tokenPrices"`
}

type alchemyTokenMeta struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Logo     string `json:"logo"`
}

type alchemyTokenPrice struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

type alchemyTransfersResponse struct {
	Result struct {
		Transfers []alchemyTransfer `json:"transfers"`
	} `json:"result"`
}

type alchemyTransfer struct {
	BlockNum string  `json:"blockNum"`
	Hash     string  `json:"hash"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	Value    float64 `json:"value"`
	Asset    string  `json:"asset"`
	Category string  `json:"category"`
	Metadata struct {
		BlockTimestamp string `json:"blockTimestamp"`
	} `json:"metadata"`
}

type alchemyBalanceResponse struct {
	Result string `json:"result"`
}

type WalletSnapshot struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_address" json:"user_id"`
	Address     string    `gorm:"not null;uniqueIndex:idx_user_address" json:"address"`
	BalanceWei  string    `gorm:"not null" json:"balance_wei"`
	BalanceEth  string    `gorm:"not null" json:"balance_eth"`
	RawResponse string    `gorm:"type:jsonb" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Transaction struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Asset       string `json:"asset"`
	Timestamp   string `json:"timestamp"`
	BlockNumber string `json:"block_number"`
	Direction   string `json:"direction"`
}

type TokenHolding struct {
	Contract string `json:"contract"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Balance  string `json:"balance"`
	USDPrice string `json:"usd_price"`
	USDValue string `json:"usd_value"`
	Logo     string `json:"logo"`
}

type Portfolio struct {
	TotalUSD string         `json:"total_usd"`
	Tokens   []TokenHolding `json:"tokens"`
}

type Counterparty struct {
	Address string `json:"address"`
	Count   int    `json:"count"`
	Label   string `json:"label,omitempty"`
}

type Overview struct {
	Address           string         `json:"address"`
	TxCount           int            `json:"tx_count"`
	IncomingCount     int            `json:"incoming_count"`
	OutgoingCount     int            `json:"outgoing_count"`
	FirstTx           string         `json:"first_tx"`
	LastTx            string         `json:"last_tx"`
	AgeDays           int            `json:"age_days"`
	TopCounterparties []Counterparty `json:"top_counterparties"`
}

type GraphNode struct {
	ID    string `json:"id"`
	Label string `json:"label,omitempty"`
	Type  string `json:"type"`
}

type GraphEdge struct {
	Source    string `json:"source"`
	Target    string `json:"target"`
	Count     int    `json:"count"`
	Direction string `json:"direction"`
}

type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type TokenCache struct {
	Address    string    `gorm:"primaryKey" json:"address"`
	TotalUSD   string    `gorm:"not null" json:"total_usd"`
	TokensJSON string    `gorm:"type:jsonb;not null" json:"-"`
	FetchedAt  time.Time `json:"fetched_at"`
}
