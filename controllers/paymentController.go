package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

// ConfirmTransaction godoc
// @Summary Example confirm a payment
// @Description example return confirms a payment on Hedera
// @Tags Payment
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /confirm-payment [get]
func ConfirmTransaction(c *gin.Context) {
	rawTxID := "0.0.6357764@1752648750.206207259"

	mirrorTxID, err := utils.ConvertToMirrorTxID(rawTxID)
	if err != nil {
		fmt.Println("❌ Error:", err)
	}

	url := fmt.Sprintf("https://testnet.mirrornode.hedera.com/api/v1/transactions/%s", mirrorTxID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(500, gin.H{"error": "Failed to query mirror node", "transactionId": mirrorTxID})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode mirror node response"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Transaction confirmed",
		"content": result,
	})
}
