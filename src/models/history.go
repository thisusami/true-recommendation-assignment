package models

import "time"

type WatchHistory struct {
	Id        string `json:"id" db:"id"`
	UserId    string `json:"user_id" db:"user_id"`
	ContentId string `json:"content_id" db:"content_id"`
	WatchedAt int64  `json:"watched_at" db:"watched_at"`
}
type WatchHistoryWithContent struct {
	Id        int64      `json:"id,omitempty" db:"id,omitempty"`
	ContentId int64      `json:"content_id,omitempty" db:"content_id,omitempty"`
	WatchedAt *time.Time `json:"watched_at,omitempty" db:"watched_at,omitempty"`
	Title     string     `json:"title,omitempty" db:"title,omitempty"`
	Genre     string     `json:"genre,omitempty" db:"genre,omitempty"`
}
