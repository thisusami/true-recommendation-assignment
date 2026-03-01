package models

import "time"

type User struct {
	Id               int64      `json:"id" db:"id"`
	Age              int16      `json:"age" db:"age"`
	Country          string     `json:"country" db:"country"`
	SubscriptionType string     `json:"subscription_type" db:"subscription_type"`
	CreatedAt        *time.Time `json:"created_at" db:"created_at"`
}
