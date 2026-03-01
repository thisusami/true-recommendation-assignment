package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	modelclient "github.com/thisusami/true-recommendation-assignment/src/model-client"
	"github.com/thisusami/true-recommendation-assignment/src/models"
	"github.com/thisusami/true-recommendation-assignment/src/repositories"
	"github.com/thisusami/true-recommendation-assignment/src/util"
	"golang.org/x/sync/semaphore"
)

type Service struct {
	DB    *repositories.Database
	Cache *repositories.Caching
	Ml    *modelclient.ModelClient
	Sem   *semaphore.Weighted
}

func (s *Service) GetRecommendationsByUser(userId string, limit int) (int, *fiber.Map) {
	key := fmt.Sprintf("rec:user:%s:limit:%d", userId, limit)
	cacheHit := false
	if cached, err := s.Cache.Get(key); err == nil {
		var result models.GetbyIdResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			cacheHit = true
			result.Metadata.CacheHit = &cacheHit
			return fiber.StatusOK, util.ToFiberMap(result)
		}

	}
	user, err := s.DB.GetUserById(userId)
	if err != nil {
		code, jsonData := models.InternalServerError.ToFiberMap()
		return code, jsonData
	}
	if user == nil {
		code, jsonData := models.NotFoundError.Set(fmt.Sprintf("User with ID %s does not exist", userId)).ToFiberMap()
		return code, jsonData
	}
	recommendationItem, err := s.Ml.GenerateRecommendations(*user, limit)
	if err != nil {
		code, jsonData := models.ServiceUnavailableError.ToFiberMap()
		return code, jsonData
	}
	var (
		recommendation []models.RecommendationItem
		now            = time.Now()
	)
	util.MaptoStruct[[]map[string]any, []models.RecommendationItem](recommendationItem, &recommendation)
	result := models.GetbyIdResponse{
		UserId:          user.Id,
		Recommendations: recommendation,
		Metadata: models.Metadata{
			CacheHit:    &cacheHit,
			GeneratedAt: &now,
			TotalCount:  len(recommendation),
		},
	}
	response := util.ToFiberMap(result)
	s.Cache.Set(key, util.ToJsonString(result), 10*time.Minute)

	return fiber.StatusOK, response
}
func (s *Service) GetBatchRecommendations(page, limit int) (int, *fiber.Map) {
	var (
		users     []models.User
		err       error
		startTime = time.Now()
	)

	users, err = s.DB.GetAllUsers(page, limit)
	if err != nil {
		code, jsonData := models.InternalServerError.ToFiberMap()
		return code, jsonData
	}
	if len(users) == 0 {
		code, jsonData := models.NotFoundError.Set("No users found for the given page and limit").ToFiberMap()
		return code, jsonData
	}
	var (
		wg sync.WaitGroup
		ch = make(chan models.Result, len(users))
	)
	for _, user := range users {
		s.Sem.Acquire(context.Background(), 1)
		wg.Add(1)
		go func(u models.User) {
			defer wg.Done()
			defer s.Sem.Release(1)
			recommendedItem, err := s.Ml.GenerateRecommendations(u, 2)
			if err != nil {
				ch <- models.Result{
					UserId: u.Id,
					Status: "failed",
					Error:  err.Error(),
				}
				return
			}
			var recommendation []models.RecommendationItem
			util.MaptoStruct[[]map[string]any, []models.RecommendationItem](recommendedItem, &recommendation)
			ch <- models.Result{
				UserId:          u.Id,
				Recommendations: recommendation,
				Status:          "success",
			}
		}(user)
	}
	wg.Wait()
	close(ch)
	var (
		results      []models.Result
		successCount = 0
		failedCount  = 0
	)
	for result := range ch {
		results = append(results, result)
		if result.Status == "success" {
			successCount++
		} else {
			failedCount++
		}
	}
	now := time.Now()
	batchResponse := models.BatchResponse{
		Page:       page,
		Limit:      limit,
		TotalUsers: len(users),
		Results:    results,
		Summary: models.Summary{
			SuccessCount:   successCount,
			FailedCount:    failedCount,
			ProcessingTime: int(time.Since(startTime).Milliseconds()),
		},
		Metadata: models.Metadata{
			GeneratedAt: &now,
		},
	}
	return fiber.StatusOK, util.ToFiberMap(batchResponse)
}

func NewService(db *repositories.Database, cache *repositories.Caching, ml *modelclient.ModelClient) *Service {
	return &Service{
		DB:    db,
		Cache: cache,
		Ml:    ml,
		Sem:   semaphore.NewWeighted(50),
	}
}
