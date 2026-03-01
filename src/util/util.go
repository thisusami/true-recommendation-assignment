package util

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func MaptoStruct[T any, U any](m T, result *U) {
	jsonData, _ := json.Marshal(m)
	_ = json.Unmarshal(jsonData, result)
}
func ToJsonString(data any) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
func ToFiberMap(data any) *fiber.Map {
	var c fiber.Map
	jsonData, _ := json.Marshal(data)
	_ = json.Unmarshal(jsonData, &c)
	return &c
}
