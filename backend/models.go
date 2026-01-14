package main

// Review represents a single customer review
type Review struct {
	ID         string `json:"id"`
	Date       string `json:"date"`
	UserID     string `json:"user_id"`
	ReviewText string `json:"review_text"`
	Rating     int    `json:"rating"`
	Source     string `json:"source"`
}

// ReviewCollection holds a list of reviews with metadata
type ReviewCollection struct {
	Reviews []Review `json:"reviews"`
	Type    string   `json:"type"` // "pre_launch" or "post_launch"
	Count   int      `json:"count"`
}

// SentimentResult represents sentiment analysis for a review
type SentimentResult struct {
	ReviewID  string  `json:"review_id"`
	Sentiment string  `json:"sentiment"` // positive, negative, neutral
	Score     float64 `json:"score"`     // confidence score
}

// ThemeResult represents an extracted theme
type ThemeResult struct {
	Theme       string `json:"theme"`
	PreCount    int    `json:"pre_count"`
	PostCount   int    `json:"post_count"`
	ChangeRate  float64 `json:"change_rate"` // percentage change
	Sentiment   string  `json:"sentiment"`   // overall sentiment for this theme
}

// SentimentSummary aggregates sentiment data
type SentimentSummary struct {
	Positive int     `json:"positive"`
	Negative int     `json:"negative"`
	Neutral  int     `json:"neutral"`
	Average  float64 `json:"average_rating"`
}

// ComparisonResult holds the comparison between pre and post launch
type ComparisonResult struct {
	PreLaunchSentiment  SentimentSummary `json:"pre_launch_sentiment"`
	PostLaunchSentiment SentimentSummary `json:"post_launch_sentiment"`
	SentimentShift      float64          `json:"sentiment_shift"` // positive = improvement
	Themes              []ThemeResult    `json:"themes"`
}

// ImpactSummary provides the overall launch impact analysis
type ImpactSummary struct {
	OverallSuccess    bool     `json:"overall_success"`
	SuccessScore      float64  `json:"success_score"` // 0-100
	KeyImprovements   []string `json:"key_improvements"`
	CriticalIssues    []string `json:"critical_issues"`
	Recommendations   []string `json:"recommendations"`
	ExecutiveSummary  string   `json:"executive_summary"`
}

// AnalysisResult is the complete analysis response
type AnalysisResult struct {
	PreLaunchReviews  ReviewCollection  `json:"pre_launch_reviews"`
	PostLaunchReviews ReviewCollection  `json:"post_launch_reviews"`
	Comparison        ComparisonResult  `json:"comparison"`
	Impact            ImpactSummary     `json:"impact"`
	AnalyzedAt        string            `json:"analyzed_at"`
}

// UploadResponse is returned after successful file upload
type UploadResponse struct {
	Success          bool   `json:"success"`
	PreLaunchCount   int    `json:"pre_launch_count"`
	PostLaunchCount  int    `json:"post_launch_count"`
	Message          string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// HealthResponse for health check endpoint
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
