package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// ReviewParser defines the interface for parsing review data
type ReviewParser interface {
	ParseCSV(reader io.Reader) ([]Review, error)
}

// AnalysisService defines the interface for the analysis service
type AnalysisService interface {
	Analyze(preReviews, postReviews []Review) (*AnalysisResult, error)
}

// CSVReviewParser implements ReviewParser for CSV files
type CSVReviewParser struct{}

// NewCSVReviewParser creates a new CSV parser instance
func NewCSVReviewParser() *CSVReviewParser {
	return &CSVReviewParser{}
}

// ParseCSV parses CSV data into Review structs
func (p *CSVReviewParser) ParseCSV(reader io.Reader) ([]Review, error) {
	csvReader := csv.NewReader(reader)
	
	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Create column index map
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[col] = i
	}

	var reviews []Review
	lineNum := 1

	for {
		lineNum++
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading line %d: %w", lineNum, err)
		}

		review := Review{}

		// Parse each field with safe access
		if idx, ok := colIndex["id"]; ok && idx < len(record) {
			review.ID = record[idx]
		}
		if idx, ok := colIndex["date"]; ok && idx < len(record) {
			review.Date = record[idx]
		}
		if idx, ok := colIndex["user_id"]; ok && idx < len(record) {
			review.UserID = record[idx]
		}
		if idx, ok := colIndex["review_text"]; ok && idx < len(record) {
			review.ReviewText = record[idx]
		}
		if idx, ok := colIndex["rating"]; ok && idx < len(record) {
			if rating, err := strconv.Atoi(record[idx]); err == nil {
				review.Rating = rating
			}
		}
		if idx, ok := colIndex["source"]; ok && idx < len(record) {
			review.Source = record[idx]
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// DefaultAnalysisService implements AnalysisService
type DefaultAnalysisService struct {
	llmClient LLMAnalyzer
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(llmClient LLMAnalyzer) *DefaultAnalysisService {
	return &DefaultAnalysisService{
		llmClient: llmClient,
	}
}

// Analyze performs the complete analysis of pre and post launch reviews
func (s *DefaultAnalysisService) Analyze(preReviews, postReviews []Review) (*AnalysisResult, error) {
	// Create review collections
	preCollection := ReviewCollection{
		Reviews: preReviews,
		Type:    "pre_launch",
		Count:   len(preReviews),
	}
	postCollection := ReviewCollection{
		Reviews: postReviews,
		Type:    "post_launch",
		Count:   len(postReviews),
	}

	// Analyze sentiments for both collections
	preSentiments, err := s.llmClient.AnalyzeSentiments(preReviews)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze pre-launch sentiments: %w", err)
	}

	postSentiments, err := s.llmClient.AnalyzeSentiments(postReviews)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze post-launch sentiments: %w", err)
	}

	// Calculate sentiment summaries
	preSummary := calculateSentimentSummary(preSentiments, preReviews)
	postSummary := calculateSentimentSummary(postSentiments, postReviews)

	// Extract themes
	themes, err := s.llmClient.ExtractThemes(preReviews, postReviews)
	if err != nil {
		return nil, fmt.Errorf("failed to extract themes: %w", err)
	}

	// Calculate sentiment shift
	sentimentShift := calculateSentimentShift(preSummary, postSummary)

	comparison := ComparisonResult{
		PreLaunchSentiment:  preSummary,
		PostLaunchSentiment: postSummary,
		SentimentShift:      sentimentShift,
		Themes:              themes,
	}

	// Generate impact summary
	impact, err := s.llmClient.GenerateImpactSummary(preCollection, postCollection, comparison)
	if err != nil {
		return nil, fmt.Errorf("failed to generate impact summary: %w", err)
	}

	result := &AnalysisResult{
		PreLaunchReviews:  preCollection,
		PostLaunchReviews: postCollection,
		Comparison:        comparison,
		Impact:            *impact,
		AnalyzedAt:        time.Now().Format(time.RFC3339),
	}

	return result, nil
}

// calculateSentimentSummary calculates aggregated sentiment stats
func calculateSentimentSummary(sentiments []SentimentResult, reviews []Review) SentimentSummary {
	summary := SentimentSummary{}

	for _, s := range sentiments {
		switch s.Sentiment {
		case "positive":
			summary.Positive++
		case "negative":
			summary.Negative++
		default:
			summary.Neutral++
		}
	}

	// Calculate average rating
	if len(reviews) > 0 {
		total := 0
		for _, r := range reviews {
			total += r.Rating
		}
		summary.Average = float64(total) / float64(len(reviews))
	}

	return summary
}

// calculateSentimentShift calculates the percentage shift in positive sentiment
func calculateSentimentShift(pre, post SentimentSummary) float64 {
	preTotal := pre.Positive + pre.Negative + pre.Neutral
	postTotal := post.Positive + post.Negative + post.Neutral

	if preTotal == 0 || postTotal == 0 {
		return 0
	}

	prePositiveRate := float64(pre.Positive) / float64(preTotal) * 100
	postPositiveRate := float64(post.Positive) / float64(postTotal) * 100

	return postPositiveRate - prePositiveRate
}

// APIHandler handles HTTP requests
type APIHandler struct {
	parser          ReviewParser
	analysisService AnalysisService
	preReviews      []Review
	postReviews     []Review
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(parser ReviewParser, analysisService AnalysisService) *APIHandler {
	return &APIHandler{
		parser:          parser,
		analysisService: analysisService,
	}
}

// HandleHealth handles the health check endpoint
func (h *APIHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
	}
	respondJSON(w, http.StatusOK, response)
}

// HandleUpload handles CSV file uploads
func (h *APIHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	// Parse multipart form with 10MB max memory
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Failed to parse form", err.Error())
		return
	}

	// Parse pre-launch file
	preLaunchFile, _, err := r.FormFile("preLaunch")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Pre-launch file is required", err.Error())
		return
	}
	defer preLaunchFile.Close()

	h.preReviews, err = h.parser.ParseCSV(preLaunchFile)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to parse pre-launch CSV", err.Error())
		return
	}

	// Parse post-launch file
	postLaunchFile, _, err := r.FormFile("postLaunch")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Post-launch file is required", err.Error())
		return
	}
	defer postLaunchFile.Close()

	h.postReviews, err = h.parser.ParseCSV(postLaunchFile)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to parse post-launch CSV", err.Error())
		return
	}

	response := UploadResponse{
		Success:         true,
		PreLaunchCount:  len(h.preReviews),
		PostLaunchCount: len(h.postReviews),
		Message:         "Files uploaded successfully. Ready for analysis.",
	}
	respondJSON(w, http.StatusOK, response)
}

// HandleAnalyze handles the analysis request
func (h *APIHandler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	if len(h.preReviews) == 0 || len(h.postReviews) == 0 {
		respondError(w, http.StatusBadRequest, "Please upload CSV files first", "")
		return
	}

	result, err := h.analysisService.Analyze(h.preReviews, h.postReviews)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Analysis failed", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// Helper functions for HTTP responses

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message, details string) {
	response := ErrorResponse{
		Error:   message,
		Code:    status,
		Details: details,
	}
	respondJSON(w, status, response)
}
