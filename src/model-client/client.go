package modelclient

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/thisusami/true-recommendation-assignment/src/models"
	"github.com/thisusami/true-recommendation-assignment/src/repositories"
	"github.com/thisusami/true-recommendation-assignment/src/util"
)

type ModelClient struct {
	Db *repositories.Database
}

func (mc *ModelClient) GenerateRecommendations(user models.User, limit int) ([]map[string]any, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()
	watch_history, err := mc.Db.GetWatchHistoryJoinContentByUserId(user.Id, 1, 50)
	if err != nil {
		return nil, err
	}
	watchedContentIds := make([]string, len(watch_history))
	for _, item := range watch_history {
		watchedContentIds = append(watchedContentIds, fmt.Sprintf("%d", item.Id))
	}
	genreFrequency := countGenresFrequency(watch_history)
	contentCandidates, err := mc.Db.GetCandidateContents(user.Id, user.Country, user.SubscriptionType, 100)
	if err != nil {
		return nil, err
	}
	latency := rand.Intn(50-30+1) + 30
	time.Sleep(time.Duration(latency) * time.Millisecond)
	if rand.Float64()*100 < 2 {
		return nil, fmt.Errorf("models interface failed")
	}
	scoredContent := make([]map[string]any, 0)
	for i := range contentCandidates {

		scoredContent = append(scoredContent, map[string]any{
			"content_id":       contentCandidates[i].Id,
			"title":            contentCandidates[i].Title,
			"genre":            contentCandidates[i].Genre,
			"popularity_score": contentCandidates[i].PopularityScore,
			"score":            scoreContent(contentCandidates[i], genreFrequency),
		})

	}

	sort.Slice(scoredContent, func(i, j int) bool {
		return scoredContent[i]["score"].(float64) > scoredContent[j]["score"].(float64)
	})
	topN := limit
	if len(scoredContent) < topN {
		topN = len(scoredContent)
	}

	// select {
	// case <-ctx.Done():
	// 	return nil, ctx.Err()
	// default:
	util.Debug("Summarize score", util.ToJsonString(scoredContent[:topN]))
	// }
	return scoredContent[:topN], nil
}
func countGenresFrequency(watch_history []models.WatchHistoryWithContent) map[string]float64 {
	genreCount := make(map[string]int)
	for _, item := range watch_history {
		genreCount[item.Genre]++
	}
	var (
		sum      = 0
		maxCount = 0
	)
	for genre, count := range genreCount {
		if count == 0 {
			delete(genreCount, genre)
		}
		sum += count
		if count > maxCount {
			maxCount = count
		}
	}
	genreFrequency := make(map[string]float64)
	if maxCount == 0 {
		return genreFrequency
	}
	for genre, count := range genreCount {
		result := float64(count) / float64(sum)
		genreFrequency[genre] = math.Trunc(result*100) / 100
	}
	util.Debug("Genre Frequency", util.ToJsonString(genreFrequency))
	return genreFrequency
}
func recencyFactor(content models.Content) float64 {
	currentTime := time.Now()
	var daysSinceCreation float64
	if content.CreatedAt != nil {
		daysSinceCreation = currentTime.Sub(*content.CreatedAt).Hours() / 24.0
	}
	recencyFactor := 1.0 / (1.0 + daysSinceCreation/365.0)
	return recencyFactor
}
func scoreContent(content models.Content, genreFrequency map[string]float64) float64 {
	popularityComponent := content.PopularityScore * 0.4

	genrePreference, ok := genreFrequency[content.Genre]
	if !ok {
		genrePreference = 0.1
	}
	genreBoost := genrePreference * 0.35

	recencyComponent := recencyFactor(content) * 0.15

	randomNoise := (rand.Float64()*0.1 - 0.05) * 0.1
	util.Debug("Score Components", fmt.Sprintf("Content ID: %d, Popularity: %.2f, Genre Boost: %.2f, Recency: %.2f, Noise: %.4f", content.Id, popularityComponent, genreBoost, recencyComponent, randomNoise))
	finalScore := popularityComponent + genreBoost + recencyComponent + randomNoise

	return math.Trunc(finalScore*100) / 100
}
func NewModelClient(db *repositories.Database) *ModelClient {
	return &ModelClient{
		Db: db,
	}
}
