package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/thisusami/true-recommendation-assignment/src/db"
	"github.com/thisusami/true-recommendation-assignment/src/handler"
	"github.com/thisusami/true-recommendation-assignment/src/models"
	"github.com/thisusami/true-recommendation-assignment/src/util"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	app := fiber.New()
	pg := db.InitPgProperty(getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/recommendation?sslmode=disable"))
	redis := db.InitRedisProperty(
		getEnv("REDIS_ADDR", "localhost:6379"),
		getEnv("REDIS_PASSWORD", ""),
	)
	defer db.Close(pg)
	defer redis.Close()
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			util.Error("Panic Recovered", e.(string))
			code, message := models.InternalServerError.ToFiberMap()
			c.Status(code).JSON(message)
		},
	}))
	handler.NewHandler(app, pg, redis).RegisterRoutes()
	app.Listen(":" + getEnv("PORT", "8080"))
}
