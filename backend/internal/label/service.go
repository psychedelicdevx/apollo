package label

import "strings"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) LabelFor(address string) string {
	l, err := s.repo.FindByAddress(strings.ToLower(address))
	if err != nil {
		return ""
	}
	return l.Name
}

func (s *Service) Get(address string) (*Label, error) {
	return s.repo.FindByAddress(strings.ToLower(address))
}

func (s *Service) Seed() error {
	n, err := s.repo.Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	return s.repo.CreateMany(seedLabels)
}

var seedLabels = []Label{
	{Address: "0x7a250d5630b4cf539739df2c5dacb4c659f2488d", Name: "Uniswap V2: Router", Category: "DEX", Source: "seed"},
	{Address: "0xe592427a0aece92de3edee1f18e0157c05861564", Name: "Uniswap V3: Router", Category: "DEX", Source: "seed"},
	{Address: "0x68b3465833fb72a70ecdf485e0e4c7bd8665fc45", Name: "Uniswap: Universal Router", Category: "DEX", Source: "seed"},
	{Address: "0xdac17f958d2ee523a2206206994597c13d831ec7", Name: "Tether: USDT", Category: "Token", Source: "seed"},
	{Address: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", Name: "Circle: USDC", Category: "Token", Source: "seed"},
	{Address: "0x6b175474e89094c44da98b954eedeac495271d0f", Name: "Maker: DAI", Category: "Token", Source: "seed"},
	{Address: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", Name: "Wrapped Ether: WETH", Category: "Token", Source: "seed"},
	{Address: "0x28c6c06298d514db089934071355e5743bf21d60", Name: "Binance 14", Category: "Exchange", Source: "seed"},
	{Address: "0x21a31ee1afc51d94c2efccaa2092ad1028285549", Name: "Binance 15", Category: "Exchange", Source: "seed"},
	{Address: "0xdfd5293d8e347dfe59e90efd55b2956a1343963d", Name: "Binance 16", Category: "Exchange", Source: "seed"},
	{Address: "0x00000000219ab540356cbb839cbe05303d7705fa", Name: "ETH2: Deposit Contract", Category: "Contract", Source: "seed"},
	{Address: "0x00000000006c3852cbef3e08e8df289169ede581", Name: "OpenSea: Seaport 1.1", Category: "NFT", Source: "seed"},
}
