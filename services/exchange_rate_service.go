package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchHBARRate() (float64, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=hedera-hashgraph&vs_currencies=usd"

	res, err := http.Get(url)
	if err != nil {
		return 0, nil
	}

	defer res.Body.Close()

	var result map[string]map[string]float64

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, nil
	}

	usd := result["hedera-hashgraph"]["usd"]
	if usd == 0 {
		return 0, fmt.Errorf("invalid exchange rate")
	}

	// rate -> 1 USD = x HBAR
	return 1 / usd, nil
}
