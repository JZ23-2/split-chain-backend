package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/JZ23-2/splitbill-backend/dtos"
)

func SendToGemini(file io.Reader) (*dtos.ReceiptResponse, error) {
	promptBytes, err := os.ReadFile("prompts/extract_text.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt: %w", err)
	}
	prompt := string(promptBytes)

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	base64Image := base64.StdEncoding.EncodeToString(imgBytes)

	payload := dtos.GeminiRequest{
		Contents: []dtos.Content{
			{
				Parts: []dtos.Part{
					{Text: prompt},
					{
						InlineData: &dtos.InlineData{
							MimeType: "image/png",
							Data:     base64Image,
						},
					},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)

	apiURL := os.Getenv("GEMINI_API_URL") + "?key=" + os.Getenv("GEMINI_API_KEY")

	res, err := http.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		resBody, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status %d: %s", res.StatusCode, resBody)
	}

	var geminiRes dtos.GeminiResponse
	if err := json.NewDecoder(res.Body).Decode(&geminiRes); err != nil {
		return nil, fmt.Errorf("failed to decode gemini response: %v", err)
	}

	if len(geminiRes.Candidates) == 0 || len(geminiRes.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("invalid response structure")
	}

	jsonText := geminiRes.Candidates[0].Content.Parts[0].Text

	jsonText = strings.TrimSpace(jsonText)
	if strings.HasPrefix(jsonText, "```") {
		jsonText = strings.TrimPrefix(jsonText, "```json")
		jsonText = strings.TrimPrefix(jsonText, "```")
		jsonText = strings.TrimSuffix(jsonText, "```")
		jsonText = strings.TrimSpace(jsonText)
	}

	var receipt dtos.ReceiptResponse
	if err := json.Unmarshal([]byte(jsonText), &receipt); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from gemini response: %v", err)
	}

	return &receipt, nil
}
