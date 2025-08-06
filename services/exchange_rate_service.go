package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func FetchHBARRate() (float64, error) {
	url := "https://api.binance.com/api/v3/ticker/price?symbol=HBARUSDT"

	res, err := http.Get(url)
	if err != nil {
		return 0, nil
	}

	defer res.Body.Close()

	var result struct {
		Price string `json:"price"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, nil
	}

	usd, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid exchange rate")
	}

	// rate -> 1 USD = x HBAR
	return 1 / usd, nil
}
