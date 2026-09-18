package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ai-translator/internal/store"
	"ai-translator/internal/translate"
)

func SetupRoutes(r *gin.Engine, s *store.Store, tc *translate.TranslatorClient) {
	// Translate endpoint
	r.POST("/api/translate", func(c *gin.Context) {
		var req translate.TranslateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := tc.Translate(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Save to history
		recordID := uuid.New().String()
		s.SaveTranslation(&store.TranslationRecord{
			ID:           recordID,
			Text:         req.Text,
			SourceLang:   req.Source,
			TargetLang:   req.Target,
			Translated:   result.TranslatedText,
			CreatedAt:    time.Now(),
		})

		c.JSON(http.StatusOK, gin.H{
			"id":              recordID,
			"translated_text": result.TranslatedText,
			"model":           result.Model,
			"tokens":          result.Tokens,
		})
	})

	// Batch translate endpoint
	r.POST("/api/translate/batch", func(c *gin.Context) {
		var req translate.BatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		results := make([]translate.TranslateResult, 0, len(req.Items))
		failed := 0

		for _, item := range req.Items {
			// Use batch source/target if not specified
			if item.Source == "" {
				item.Source = req.Source
			}
			if item.Target == "" {
				item.Target = req.Target
			}

			result, err := tc.Translate(item)
			if err != nil {
				failed++
				results = append(results, translate.TranslateResult{
					TranslatedText: fmt.Sprintf("[Error: %v]", err),
					Model:          "error",
					Tokens:         0,
				})
				continue
			}

			results = append(results, *result)

			// Save each translation
			recordID := uuid.New().String()
			s.SaveTranslation(&store.TranslationRecord{
				ID:           recordID,
				Text:         item.Text,
				SourceLang:   item.Source,
				TargetLang:   item.Target,
				Translated:   result.TranslatedText,
				CreatedAt:    time.Now(),
			})
		}

		c.JSON(http.StatusOK, translate.BatchResult{
			Results: results,
			Failed:  failed,
			Total:   len(req.Items),
		})
	})

	// Detect language endpoint
	r.POST("/api/detect", func(c *gin.Context) {
		var req struct {
			Text string `json:"text" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		detected := translate.DetectLanguage(req.Text)
		c.JSON(http.StatusOK, detected)
	})

	// Get history endpoint
	r.GET("/api/history", func(c *gin.Context) {
		limit := 20
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}
		records, err := s.GetTranslations(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, records)
	})

	// Get languages endpoint
	r.GET("/api/languages", func(c *gin.Context) {
		c.JSON(http.StatusOK, translate.Languages)
	})

	// Save term endpoint
	r.POST("/api/terms", func(c *gin.Context) {
		var req struct {
			Source   string `json:"source"`
			Target   string `json:"target"`
			LangPair string `json:"lang_pair"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		termID := uuid.New().String()
		s.SaveTerm(&store.Term{
			ID:        termID,
			Source:    req.Source,
			Target:    req.Target,
			LangPair:  req.LangPair,
			CreatedAt: time.Now(),
		})

		c.JSON(http.StatusOK, gin.H{"id": termID})
	})

	// Get terms endpoint
	r.GET("/api/terms", func(c *gin.Context) {
		langPair := c.Query("lang_pair")
		terms, err := s.GetTerms(langPair)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, terms)
	})

	// Delete term endpoint
	r.DELETE("/api/terms/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := s.DeleteTerm(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	})

	// Stats endpoint
	r.GET("/api/stats", func(c *gin.Context) {
		stats, err := s.GetStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)
	})
}
