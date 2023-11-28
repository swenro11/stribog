package entity

type TokenInfo struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Type        string `json:"type"`
	Decimals    int    `json:"decimals"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Explorer    string `json:"explorer"`
	Status      string `json:"status"`
	Id          string `json:"id"`
}

/* example from https://github.com/trustwallet/assets/blob/master/blockchains/optimism/assets/0x4200000000000000000000000000000000000006/info.json
{
    "name": "Wrapped Ether",
    "symbol": "WETH",
    "type": "OPTIMISM",
    "decimals": 18,
    "description": "wETH is wrapped ETH",
    "website": "https://weth.io/",
    "explorer": "https://optimistic.etherscan.io/token/0x4200000000000000000000000000000000000006",
    "status": "active",
    "id": "0x4200000000000000000000000000000000000006",
    "links": [
        {
            "name": "coinmarketcap",
            "url": "https://coinmarketcap.com/currencies/weth/"
        },
        {
            "name": "coingecko",
            "url": "https://coingecko.com/coins/weth/"
        }
    ]
}
*/
