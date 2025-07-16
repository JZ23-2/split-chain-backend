package utils

import (
	"fmt"
	"strings"
)

func ConvertToMirrorTxID(sdkTxID string) (string, error) {
	parts := strings.Split(sdkTxID, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid transaction ID format")
	}

	accountID := parts[0]
	timestamp := strings.Replace(parts[1], ".", "-", 1)

	return fmt.Sprintf("%s-%s", accountID, timestamp), nil
}
