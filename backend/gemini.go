package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LLMAnalyzer defines the interface for LLM-based analysis
type LLMAnalyzer interface {
	AnalyzeSentiments(reviews []Review) ([]SentimentResult, error)
	ExtractThemes(preReviews, postReviews []Review) ([]ThemeResult, error)
	GenerateImpactSummary(pre, post ReviewCollection, comparison ComparisonResult) (*ImpactSummary, error)
}

// GroqClient implements LLMAnalyzer using Groq's API
type GroqClient struct {
	apiKey  string
	baseURL string
	model   string
}

// NewGroqClient creates a new Groq client instance
func NewGroqClient(apiKey string) *GroqClient {
	return &GroqClient{
		apiKey:  apiKey,
		baseURL: "https://api.groq.com/openai/v1",
		model:   "llama-3.3-70b-versatile",
	}
}

// GroqRequest represents the request payload for Groq API
type GroqRequest struct {
	Model    string        `json:"model"`
	Messages []GroqMessage `json:"messages"`
}

// GroqMessage represents a message in Groq request
type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GroqResponse represents the response from Groq API
type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// callGroqAPI makes a request to the Groq API
func (g *GroqClient) callGroqAPI(prompt string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", g.baseURL)

	request := GroqRequest{
		Model: g.model,
		Messages: []GroqMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Groq API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Groq API error (status %d): %s", resp.StatusCode, string(body))
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if groqResp.Error != nil {
		return "", fmt.Errorf("Groq API error: %s", groqResp.Error.Message)
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from Groq")
	}

	return groqResp.Choices[0].Message.Content, nil
}

// AnalyzeSentiments analyzes sentiment for each review
func (g *GroqClient) AnalyzeSentiments(reviews []Review) ([]SentimentResult, error) {
	if len(reviews) == 0 {
		return []SentimentResult{}, nil
	}

	// Prepare reviews text for analysis
	reviewsText := ""
	for _, r := range reviews {
		reviewsText += fmt.Sprintf("ID: %s | Rating: %d | Review: %s\n", r.ID, r.Rating, r.ReviewText)
	}

	prompt := fmt.Sprintf(`Analyze the sentiment of these customer reviews. For each review, classify as "positive", "negative", or "neutral" with a confidence score (0-1).

Reviews:
%s

Respond ONLY with a valid JSON array in this exact format (no markdown, no explanation):
[{"review_id": "id", "sentiment": "positive/negative/neutral", "score": 0.95}]`, reviewsText)

	response, err := g.callGroqAPI(prompt)
	if err != nil {
		return nil, err
	}

	// Clean response (remove markdown code blocks if present)
	response = cleanJSONResponse(response)

	var results []SentimentResult
	if err := json.Unmarshal([]byte(response), &results); err != nil {
		return nil, fmt.Errorf("failed to parse sentiment results: %w, response: %s", err, response)
	}

	return results, nil
}

// ExtractThemes extracts and compares themes between pre and post launch reviews
func (g *GroqClient) ExtractThemes(preReviews, postReviews []Review) ([]ThemeResult, error) {
	preText := formatReviewsForThemes(preReviews)
	postText := formatReviewsForThemes(postReviews)

	prompt := fmt.Sprintf(`Analyze and compare themes between pre-launch and post-launch customer reviews.

PRE-LAUNCH REVIEWS:
%s

POST-LAUNCH REVIEWS:
%s

Extract the top 8 themes mentioned across both sets. For each theme, count occurrences in pre and post launch, calculate percentage change, and determine overall sentiment.

Respond ONLY with a valid JSON array in this exact format (no markdown, no explanation):
[{"theme": "theme name", "pre_count": 5, "post_count": 8, "change_rate": 60.0, "sentiment": "positive/negative/neutral"}]`, preText, postText)

	response, err := g.callGroqAPI(prompt)
	if err != nil {
		return nil, err
	}

	response = cleanJSONResponse(response)

	var results []ThemeResult
	if err := json.Unmarshal([]byte(response), &results); err != nil {
		return nil, fmt.Errorf("failed to parse theme results: %w, response: %s", err, response)
	}

	return results, nil
}

// GenerateImpactSummary generates an executive summary of the launch impact
func (g *GroqClient) GenerateImpactSummary(pre, post ReviewCollection, comparison ComparisonResult) (*ImpactSummary, error) {
	prompt := fmt.Sprintf(`You are analyzing the impact of a feature launch based on customer reviews.

PRE-LAUNCH DATA:
- Total reviews: %d
- Positive: %d, Negative: %d, Neutral: %d
- Average rating: %.2f

POST-LAUNCH DATA:
- Total reviews: %d
- Positive: %d, Negative: %d, Neutral: %d
- Average rating: %.2f

SENTIMENT SHIFT: %.2f%%

KEY THEMES IDENTIFIED:
%s

Based on this data, provide a comprehensive launch impact analysis.

Respond ONLY with a valid JSON object in this exact format (no markdown, no explanation):
{
  "overall_success": true/false,
  "success_score": 75.5,
  "key_improvements": ["improvement 1", "improvement 2"],
  "critical_issues": ["issue 1", "issue 2"],
  "recommendations": ["recommendation 1", "recommendation 2"],
  "executive_summary": "A 2-3 sentence summary of the launch impact"
}`,
		pre.Count,
		comparison.PreLaunchSentiment.Positive,
		comparison.PreLaunchSentiment.Negative,
		comparison.PreLaunchSentiment.Neutral,
		comparison.PreLaunchSentiment.Average,
		post.Count,
		comparison.PostLaunchSentiment.Positive,
		comparison.PostLaunchSentiment.Negative,
		comparison.PostLaunchSentiment.Neutral,
		comparison.PostLaunchSentiment.Average,
		comparison.SentimentShift,
		formatThemesForSummary(comparison.Themes))

	response, err := g.callGroqAPI(prompt)
	if err != nil {
		return nil, err
	}

	response = cleanJSONResponse(response)

	var result ImpactSummary
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse impact summary: %w, response: %s", err, response)
	}

	return &result, nil
}

// Helper functions

func cleanJSONResponse(response string) string {
	// Remove markdown code blocks if present
	response = removePrefix(response, "```json")
	response = removePrefix(response, "```")
	response = removeSuffix(response, "```")
	return response
}

func removePrefix(s, prefix string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\n' || s[0] == '\t') {
		s = s[1:]
	}
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

func removeSuffix(s, suffix string) string {
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\n' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func formatReviewsForThemes(reviews []Review) string {
	result := ""
	for _, r := range reviews {
		result += fmt.Sprintf("- %s (Rating: %d)\n", r.ReviewText, r.Rating)
	}
	return result
}

func formatThemesForSummary(themes []ThemeResult) string {
	result := ""
	for _, t := range themes {
		result += fmt.Sprintf("- %s: Pre=%d, Post=%d, Change=%.1f%%, Sentiment=%s\n",
			t.Theme, t.PreCount, t.PostCount, t.ChangeRate, t.Sentiment)
	}
	return result
}
