package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type TranslateRequest struct {
	Text     string `json:"text"`
	Source   string `json:"source"`
	Target   string `json:"target"`
	Model    string `json:"model"`
}

type TranslateResponse struct {
	TranslatedText string `json:"translated_text"`
	Model          string `json:"model"`
	Tokens         int    `json:"tokens"`
}

type TranslatorClient struct {
	APIKey  string
	BaseURL string
	Model   string
}

var Languages = []Language{
	{Code: "auto", Name: "Auto Detect"},
	{Code: "en", Name: "English"},
	{Code: "zh", Name: "Chinese (Simplified)"},
	{Code: "zh-TW", Name: "Chinese (Traditional)"},
	{Code: "ja", Name: "Japanese"},
	{Code: "ko", Name: "Korean"},
	{Code: "fr", Name: "French"},
	{Code: "de", Name: "German"},
	{Code: "es", Name: "Spanish"},
	{Code: "pt", Name: "Portuguese"},
	{Code: "ru", Name: "Russian"},
	{Code: "ar", Name: "Arabic"},
	{Code: "hi", Name: "Hindi"},
	{Code: "th", Name: "Thai"},
	{Code: "vi", Name: "Vietnamese"},
	{Code: "id", Name: "Indonesian"},
	{Code: "ms", Name: "Malay"},
	{Code: "tl", Name: "Filipino"},
	{Code: "it", Name: "Italian"},
	{Code: "nl", Name: "Dutch"},
}

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func NewClient() *TranslatorClient {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	baseURL := os.Getenv("DEEPSEEK_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}
	return &TranslatorClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   "deepseek-chat",
	}
}

func (c *TranslatorClient) Translate(req TranslateRequest) (*TranslateResponse, error) {
	if c.APIKey == "" {
		return c.mockTranslate(req), nil
	}

	prompt := fmt.Sprintf("Translate the following text from %s to %s. Return only the translation, nothing else:\n\n%s", req.Source, req.Target, req.Text)

	data := map[string]interface{}{
		"model": c.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
	}

	jsonData, _ := json.Marshal(data)

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var apiResp map[string]interface{}
	json.Unmarshal(body, &apiResp)

	choices, ok := apiResp["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("invalid response")
	}

	msg, _ := json.Marshal(choices[0].(map[string]interface{})["message"])
	var message map[string]string
	json.Unmarshal(msg, &message)

	translatedText := message["content"]
	tokens := 0
	if usage, ok := apiResp["usage"].(map[string]interface{}); ok {
		if t, ok := usage["total_tokens"].(float64); ok {
			tokens = int(t)
		}
	}

	return &TranslateResponse{
		TranslatedText: translatedText,
		Model:          c.Model,
		Tokens:         tokens,
	}, nil
}

func (c *TranslatorClient) mockTranslate(req TranslateRequest) *TranslateResponse {
	return &TranslateResponse{
		TranslatedText: "[MOCK] " + req.Text + " → " + req.Target,
		Model:          "mock",
		Tokens:         0,
	}
}
