package utils

import "fmt"

func FormatUSD(amount int) string {
	return fmt.Sprintf("$%.2f", float64(amount)/100.0)
}

func FormatUSDtoInt(amount float64) int {
	return int(amount * 100)
}
