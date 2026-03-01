package handler

import (
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/gofiber/fiber/v2"
	modelclient "github.com/thisusami/true-recommendation-assignment/src/model-client"
	"github.com/thisusami/true-recommendation-assignment/src/models"
	"github.com/thisusami/true-recommendation-assignment/src/repositories"
	"github.com/thisusami/true-recommendation-assignment/src/services"
	"github.com/thisusami/true-recommendation-assignment/src/util"
	"gorm.io/gorm"
)

type Handler struct {
	Fiber *fiber.App
	Pg    *gorm.DB
	Redis *redis.Client
}

func (h *Handler) RegisterRoutes() *Handler {
	h.Fiber.Get("/health", func(c *fiber.Ctx) error {
		util.InBound("Health Check", c.BaseURL()+""+c.Path())
		return c.JSON(&fiber.Map{"status": "ok"})
	})
	h.Fiber.Get("/users/:user_id/recommendations", h.GetRecommendations)
	h.Fiber.Get("/recommendations/batch", h.GetBatchRecommendations)
	util.Info("Routes registered successfully", "")
	return h
}
func (h *Handler) GetRecommendations(c *fiber.Ctx) error {
	util.InBound("Request Get Recommendations", c.BaseURL()+""+c.Path())
	startTime := time.Now()
	if !util.IsNotEmptyString(c.Params("user_id")) {
		code, jsonData := models.BadRequestError.Set("user_id is required").ToFiberMap()
		return c.Status(code).JSON(jsonData)
	}
	limit := c.QueryInt("limit", 10)
	if !util.MinMaxInteger(limit, 1, 50) {
		code, jsonData := models.BadRequestError.Set("limit must be between 1 and 50").ToFiberMap()
		return c.Status(code).JSON(jsonData)
	}
	pg := repositories.NewRecommendedRepository(h.Pg)
	redisClient := repositories.NewCaching(h.Redis)
	ml := modelclient.NewModelClient(pg)
	code, service := services.NewService(pg, redisClient, ml).GetRecommendationsByUser(c.Params("user_id"), limit)
	util.OutBound("Response Get Recommendations", util.ToJsonString(service), float32(time.Since(startTime).Milliseconds()))
	return c.Status(code).JSON(service)
}
func (h *Handler) GetBatchRecommendations(c *fiber.Ctx) error {
	util.InBound("Request Get Batch Recommendations", c.BaseURL()+""+c.Path())
	startTime := time.Now()
	var (
		page  = c.QueryInt("page", 1)
		limit = c.QueryInt("limit", 20)
	)
	if !util.MinMaxInteger(page, 1, -1) {
		code, jsonData := models.BadRequestError.Set("page must be between 1 and 100").ToFiberMap()
		return c.Status(code).JSON(jsonData)
	}
	if !util.MinMaxInteger(limit, 1, 100) {
		code, jsonData := models.BadRequestError.Set("limit must be between 1 and 50").ToFiberMap()
		return c.Status(code).JSON(jsonData)
	}
	pg := repositories.NewRecommendedRepository(h.Pg)
	redisClient := repositories.NewCaching(h.Redis)
	ml := modelclient.NewModelClient(pg)
	code, service := services.NewService(pg, redisClient, ml).GetBatchRecommendations(page, limit)
	util.OutBound("Response Get Batch Recommendations", util.ToJsonString(service), float32(time.Since(startTime).Milliseconds()))
	return c.Status(code).JSON(service)
}
func NewHandler(fiberApp *fiber.App, pg *gorm.DB, redisClient *redis.Client) *Handler {
	return &Handler{
		Fiber: fiberApp,
		Pg:    pg,
		Redis: redisClient,
	}
}
