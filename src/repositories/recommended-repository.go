package repositories

import (
	"fmt"
	"time"

	"github.com/thisusami/true-recommendation-assignment/src/models"
	"github.com/thisusami/true-recommendation-assignment/src/util"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func (d *Database) GetUserById(userId string) (*models.User, error) {
	var user models.User
	startTime := time.Now()
	util.Request("Request Get User By Id", userId)
	result := d.DB.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		util.Error("Error Get User By Id", result.Error.Error())
		return nil, result.Error
	}
	duration := time.Since(startTime)
	util.Response("Response Get User By Id", util.ToJsonString(user), float32(duration.Milliseconds()))
	return &user, nil
}
func (d *Database) GetAllUsers(page, limit int) ([]models.User, error) {
	startTime := time.Now()
	util.Request("Request Get All Users", fmt.Sprintf("page: %d, limit: %d", page, limit))
	var users []models.User
	offset := (page - 1) * limit
	result := d.DB.Table("users").Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		util.Error("Error Get All Users", result.Error.Error())
		return nil, result.Error
	}
	duration := time.Since(startTime)
	util.Response("Response Get All Users", util.ToJsonString(users), float32(duration.Milliseconds()))
	return users, nil
}
func (d *Database) GetWatchHistoryByUserId(userId string, page, limit int) ([]models.WatchHistory, error) {
	var history []models.WatchHistory
	startTime := time.Now()
	util.Request("Request Get Watch History By UserId", userId)
	offset := (page - 1) * limit
	result := d.DB.Table("user_watch_history").Where("user_id = ?", userId).Offset(offset).Limit(limit).Find(&history)
	if result.Error != nil {
		util.Error("Error Get Watch History By UserId", result.Error.Error())
		return nil, result.Error
	}
	duration := time.Since(startTime)
	util.Response("Response Get Watch History By UserId", util.ToJsonString(history), float32(duration.Milliseconds()))
	return history, nil
}
func (d *Database) GetWatchHistoryJoinContentByUserId(userId int64, page, limit int) ([]models.WatchHistoryWithContent, error) {
	var history []models.WatchHistoryWithContent
	startTime := time.Now()
	util.Request("Request Watch History Join Content By UserId", fmt.Sprintf("%d", userId))
	result := d.DB.Table("user_watch_history uwh").
		Select("c.id, c.genre, uwh.watched_at").
		Joins("join content c on uwh.content_id = c.id").
		Where("uwh.user_id = ?", userId).
		Order("uwh.watched_at DESC").Limit(limit).Scan(&history)
	if result.Error != nil {
		util.Error("Error Watch History Join Content By UserId", result.Error.Error())
		return nil, result.Error
	}
	duration := time.Since(startTime)
	util.Response("Response Watch History Join Content By UserId", util.ToJsonString(history), float32(duration.Milliseconds()))
	return history, nil
}
func (d *Database) GetCandidateContents(userId int64, country string, subTier string, limit int) ([]models.Content, error) {
	var contents []models.Content
	startTime := time.Now()
	util.Request("Request Candidate Contents By UserId", fmt.Sprintf("%d", userId))
	subQuery := d.DB.Table("user_watch_history").Select("content_id").Where("user_id = ?", userId)
	result := d.DB.Table("content").
		Select("id, title, genre, popularity_score, created_at").
		Where("id NOT IN (?) and available_countries @> ARRAY[?]::text[] and available_subscription @> ARRAY[?]::text[]", subQuery, country, subTier).
		Order("popularity_score DESC").Limit(limit).Scan(&contents)
	if result.Error != nil {
		util.Error("Error Candidate Contents By UserId", result.Error.Error())
		return nil, result.Error
	}
	duration := time.Since(startTime)
	util.Response("Response Candidate Contents By UserId", util.ToJsonString(contents), float32(duration.Milliseconds()))
	return contents, nil
}
func NewRecommendedRepository(db *gorm.DB) *Database {
	return &Database{
		DB: db,
	}
}
