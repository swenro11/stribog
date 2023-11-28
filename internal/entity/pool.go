// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Pool
type Pool struct {
	PriceCoinOne float64 `json:"price_coin_one"   example:"0.4"`
	PriceCoinTwo float64 `json:"price_coin_two"   example:"0.4"`
	CoinOneUID   string  `json:"coin_one_uid"     example:"0xB777ABf5100657E378E8a71D054a1BdE9224aa96"`
	CoinOneTitle string  `json:"coin_one_title"   example:"BUSD"`
	CoinTwoUID   string  `json:"coin_two_uid"     example:"0xB777ABf5100657E378E8a71D054a1BdE9224aa96"`
	CoinTwoTitle string  `json:"coin_two_title"   example:"BUSD"`
}
