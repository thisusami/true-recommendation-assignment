package util

import (
	"fmt"
	"time"
)

type Logger struct {
	Action       string     `json:"action"`
	Description  string     `json:"description,omitempty"`
	Message      string     `json:"message,omitempty"`
	ResponseTime float32    `json:"responseTime,omitempty"`
	Timestamp    *time.Time `json:"timestamp,omitempty"`
}

func Info(description, msg string) {
	now := time.Now()
	fmt.Println(ToJsonString(Logger{
		Action:      "INFO",
		Description: description,
		Message:     ToJsonString(msg),
		Timestamp:   &now,
	}))
}

func Error(description, msg string) {
	now := time.Now()
	fmt.Println(ToJsonString(Logger{
		Action:      "EXCEPTION",
		Description: description,
		Message:     ToJsonString(msg),
		Timestamp:   &now,
	}))
}
func Debug(description, msg string) {
	now := time.Now()
	fmt.Println(ToJsonString(Logger{
		Action:      "DEBUG",
		Description: description,
		Message:     msg,
		Timestamp:   &now,
	}))
}
func InBound(description, msg string) {
	now := time.Now()
	fmt.Println(ToJsonString(Logger{
		Action:      "INBOUND",
		Description: description,
		Message:     msg,
		Timestamp:   &now,
	}))
}

func OutBound(description, msg string, responseTime float32) {
	now := time.Now()
	fmt.Println(ToJsonString(Logger{
		Action:       "OUTBOUND",
		Description:  description,
		Message:      msg,
		ResponseTime: responseTime,
		Timestamp:    &now,
	}))
}
func Request(description, msg string) {
	now := time.Now()
	fmt.Println(ToJsonString(Logger{
		Action:      "REQUEST",
		Description: description,
		Message:     msg,
		Timestamp:   &now,
	}))
}
func Response(description, msg string, responseTime float32) {
	now := time.Now()
	fmt.Println(ToJsonString(Logger{
		Action:       "RESPONSE",
		Description:  description,
		Message:      msg,
		ResponseTime: responseTime,
		Timestamp:    &now,
	}))
}
