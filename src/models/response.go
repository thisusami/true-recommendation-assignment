package models

import "time"

type Metadata struct {
	CacheHit    *bool      `json:"cache_hit,omitempty"`
	GeneratedAt *time.Time `json:"generated_at"`
	TotalCount  int        `json:"total_count,omitempty"`
}
type RecommendationItem struct {
	ContentId       int32   `json:"content_id"`
	Title           string  `json:"title"`
	Genre           string  `json:"genre"`
	PopularityScore float64 `json:"popularity_score"`
	Score           float64 `json:"score"`
}
type GetbyIdResponse struct {
	UserId          int64                `json:"user_id"`
	Recommendations []RecommendationItem `json:"recommendations"`
	Metadata        Metadata             `json:"metadata"`
}
type Result struct {
	UserId          int64                `json:"user_id"`
	Recommendations []RecommendationItem `json:"recommendations"`
	Status          string               `json:"status"`
	Error           string               `json:"error,omitempty"`
	Message         string               `json:"message,omitempty"`
}
type Summary struct {
	SuccessCount   int `json:"success_count"`
	FailedCount    int `json:"failed_count"`
	ProcessingTime int `json:"processing_time_ms"`
}
type BatchResponse struct {
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	TotalUsers int      `json:"total_users"`
	Results    []Result `json:"results"`
	Summary    Summary  `json:"summary"`
	Metadata   Metadata `json:"metadata"`
}
