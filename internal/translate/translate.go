package translate

import (
	"fmt"
	"strings"
	"time"
)

const (
	DefaultSource = "en"
	DefaultTarget = "es"
	MaxTextLength = 5000
)

type TranslatorClient struct {
	BaseURL     string
	APIKey      string
	DefaultMode string
}

type TranslateRequest struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	Target string `json:"target"`
	Model  string `json:"model"`
}

type TranslateResult struct {
	TranslatedText string `json:"translated_text"`
	Model          string `json:"model"`
	Tokens         int    `json:"tokens"`
	Duration       int64  `json:"duration_ms"`
}

type BatchRequest struct {
	Items  []TranslateRequest `json:"items"`
	Source string             `json:"source"`
	Target string             `json:"target"`
}

type BatchResult struct {
	Results []TranslateResult `json:"results"`
	Failed  int               `json:"failed"`
	Total   int               `json:"total"`
}

type DetectedLanguage struct {
	Language   string  `json:"language"`
	Confidence float64 `json:"confidence"`
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

func NewTranslatorClient(baseURL, apiKey string) *TranslatorClient {
	return &TranslatorClient{
		BaseURL:     baseURL,
		APIKey:      apiKey,
		DefaultMode: "fast",
	}
}

func (tc *TranslatorClient) Translate(req TranslateRequest) (*TranslateResult, error) {
	if req.Source == "" {
		req.Source = DefaultSource
	}
	if req.Target == "" {
		req.Target = DefaultTarget
	}
	if len(req.Text) > MaxTextLength {
		req.Text = req.Text[:MaxTextLength]
	}

	start := time.Now()

	result := &TranslateResult{
		TranslatedText: simulateTranslation(req.Text, req.Target),
		Model:          "deepseek",
		Tokens:         len(strings.Fields(req.Text)),
		Duration:       time.Since(start).Milliseconds(),
	}

	return result, nil
}

func (tc *TranslatorClient) BatchTranslate(req BatchRequest) *BatchResult {
	results := make([]TranslateResult, 0, len(req.Items))
	failed := 0

	for i := range req.Items {
		if req.Items[i].Source == "" {
			req.Items[i].Source = req.Source
		}
		if req.Items[i].Target == "" {
			req.Items[i].Target = req.Target
		}

		result, err := tc.Translate(req.Items[i])
		if err != nil {
			failed++
			results = append(results, TranslateResult{
				TranslatedText: fmt.Sprintf("[Error: %v]", err),
				Model:          "error",
				Tokens:         0,
			})
			continue
		}
		results = append(results, *result)
	}

	return &BatchResult{
		Results: results,
		Failed:  failed,
		Total:   len(req.Items),
	}
}

func (tc *TranslatorClient) DetectLanguage(text string) *DetectedLanguage {
	return DetectLanguage(text)
}

func (tc *TranslatorClient) GetSupportedLanguages() []Language {
	return Languages
}

func (tc *TranslatorClient) HealthCheck() bool {
	return true
}

func simulateTranslation(text, targetLang string) string {
	prefix := map[string]string{
		"es": "ES:", "fr": "FR:", "de": "DE:", "zh": "ZH:",
		"ja": "JA:", "ko": "KO:", "ar": "AR:", "pt": "PT:",
		"ru": "RU:", "hi": "HI:", "it": "IT:", "nl": "NL:",
		"tr": "TR:", "vi": "VI:", "th": "TH:",
	}

	prefixStr := ""
	if p, ok := prefix[targetLang]; ok {
		prefixStr = p + " "
	}

	return prefixStr + "[" + targetLang + "] " + text
}

func DetectLanguage(text string) *DetectedLanguage {
	lower := text[:min(100, len(text))]

	hasChinese, hasJapanese, hasKorean, hasArabic, hasCyrillic := false, false, false, false, false
	for _, r := range lower {
		switch {
		case r >= 0x4E00 && r <= 0x9FFF:
			hasChinese = true
		case r >= 0x3040 && r <= 0x30FF:
			hasJapanese = true
		case r >= 0xAC00 && r <= 0xD7A3:
			hasKorean = true
		case r >= 0x0600 && r <= 0x06FF:
			hasArabic = true
		case r >= 0x0400 && r <= 0x04FF:
			hasCyrillic = true
		}
	}

	if hasJapanese {
		return &DetectedLanguage{Language: "ja", Confidence: 0.90}
	}
	if hasKorean {
		return &DetectedLanguage{Language: "ko", Confidence: 0.90}
	}
	if hasChinese {
		return &DetectedLanguage{Language: "zh", Confidence: 0.85}
	}
	if hasArabic {
		return &DetectedLanguage{Language: "ar", Confidence: 0.88}
	}
	if hasCyrillic {
		return &DetectedLanguage{Language: "ru", Confidence: 0.87}
	}

	commonWords := map[string]map[string]int{
		"en": {"the": 1, "and": 1, "is": 1, "are": 1, "was": 1},
		"es": {"el": 1, "la": 1, "los": 1, "de": 1, "que": 1},
		"fr": {"le": 1, "la": 1, "les": 1, "de": 1, "que": 1},
		"de": {"der": 1, "die": 1, "das": 1, "und": 1, "ist": 1},
		"pt": {"o": 1, "a": 1, "os": 1, "de": 1, "e": 1},
		"it": {"il": 1, "la": 1, "i": 1, "di": 1, "e": 1},
	}

	for lang, words := range commonWords {
		count := 0
		for word := range words {
			if containsWord(lower, word) {
				count++
			}
		}
		if count >= 2 {
			return &DetectedLanguage{Language: lang, Confidence: float64(count) / float64(len(words))}
		}
	}

	return &DetectedLanguage{Language: "en", Confidence: 0.50}
}

func containsWord(text, word string) bool {
	return len(text) > 0 && (len(word) == 0 || text[0:len(word)] == word ||
		(len(text) > len(word) && text[len(text)-len(word):] == word) ||
		containsSubstring(text, " "+word+" "))
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateID() string {
	return fmt.Sprintf("trans-%d", time.Now().UnixNano())
}
