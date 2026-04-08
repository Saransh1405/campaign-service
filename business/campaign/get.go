package campaign

import (
	"campaign-service/constants"
	"campaign-service/logger"
	"campaign-service/models"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"campaign-service/library/postgres"
	"campaign-service/library/redis_provider"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetCampaign(ctx *gin.Context, request *models.GetCampaignRequest) (interface{}, int64, error) {
	log := logger.GetLoggerWithoutContext()

	userID := request.UserID
	if userID == "" {
		log.With(zap.Error(errors.New(constants.UserNotFoundMessage))).Error(constants.UserNotFoundMessage)
		return nil, 0, errors.New(constants.UserNotFoundMessage)
	}

	// user, err := helperfunctions.ValidateUserExists(ctx, userID)
	// if err != nil {
	// 	log.With(zap.Error(err)).Error(constants.UserNotFoundMessage)
	// 	return nil, 0, err
	// }

	// if !user.EmailVerified {
	// 	log.With(zap.Error(errors.New(constants.UserNotVerifiedMessage))).Error(constants.UserNotVerifiedMessage)
	// 	return nil, 0, errors.New(constants.UserNotVerifiedMessage)
	// }

	redis := redis_provider.Client
	cacheKey := fmt.Sprintf("campaign:user:%s", userID)
	if request.ID != "" {
		cacheKey += fmt.Sprintf(":%s", request.ID)
	}
	if request.City != "" {
		cacheKey += fmt.Sprintf(":city:%s", request.City)
	}
	if request.State != "" {
		cacheKey += fmt.Sprintf(":state:%s", request.State)
	}
	if request.Country != "" {
		cacheKey += fmt.Sprintf(":country:%s", request.Country)
	}
	if request.MinPrice != 0 {
		cacheKey += fmt.Sprintf(":min_price:%d", request.MinPrice)
	}
	if request.MaxPrice != 0 {
		cacheKey += fmt.Sprintf(":max_price:%d", request.MaxPrice)
	}
	if request.StartDate != "" {
		cacheKey += fmt.Sprintf(":start_date:%s", request.StartDate)
	}
	if request.EndDate != "" {
		cacheKey += fmt.Sprintf(":end_date:%s", request.EndDate)
	}
	if request.Status != "" {
		cacheKey += fmt.Sprintf(":status:%s", request.Status)
	}
	if len(request.Tags) > 0 {
		cacheKey += fmt.Sprintf(":tags:%s", strings.Join(request.Tags, ","))
	}
	if request.SortBy != "" {
		cacheKey += fmt.Sprintf(":sort_by:%s", request.SortBy)
	}
	if request.SortOrder != "" {
		cacheKey += fmt.Sprintf(":sort_order:%s", request.SortOrder)
	}
	if request.Skip > 0 {
		cacheKey += fmt.Sprintf(":skip:%d", request.Skip)
	}
	if request.Limit > 0 {
		cacheKey += fmt.Sprintf(":limit:%d", request.Limit)
	}

	cachedResult, err := redis.Get(ctx, cacheKey).Result()
	if err == nil && cachedResult != "" {
		var cachedResponse struct {
			Campaigns []models.Campaign `json:"campaigns"`
			Total     int64             `json:"total"`
		}
		if err := json.Unmarshal([]byte(cachedResult), &cachedResponse); err == nil {
			return cachedResponse.Campaigns, cachedResponse.Total, nil
		}
	}

	db := postgres.DB
	query := db.Model(&models.Campaign{}).Preload("Location").Preload("Participants").Preload("StatusLogs").Where("user_id = ?", userID)

	if request.ID != "" {
		query = query.Where("id = ?", request.ID)
	}
	if request.City != "" {
		query = query.Where("city = ?", request.City)
	}
	if request.State != "" {
		query = query.Where("state = ?", request.State)
	}
	if request.Country != "" {
		query = query.Where("country = ?", request.Country)
	}
	if request.MinPrice != 0 {
		query = query.Where("price >= ?", request.MinPrice)
	}
	if request.MaxPrice != 0 {
		query = query.Where("price <= ?", request.MaxPrice)
	}
	if request.StartDate != "" {
		query = query.Where("start_date = ?", request.StartDate)
	}
	if request.EndDate != "" {
		query = query.Where("end_date = ?", request.EndDate)
	}
	if request.Status != "" {
		query = query.Where("status = ?", request.Status)
	}

	if len(request.Tags) > 0 {
		query = query.Where("tags ?| ?", request.Tags)
	}

	if request.SortBy != "" {
		sortOrder := "ASC"
		if request.SortOrder != "" {
			sortOrder = strings.ToUpper(request.SortOrder)
		}
		query = query.Order(fmt.Sprintf("%s %s", request.SortBy, sortOrder))
	}

	var campaigns []models.Campaign
	var total int64

	if err := query.Count(&total).Error; err != nil {
		log.Error("Failed to count campaigns", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count campaigns: %w", err)
	}

	if request.Skip > 0 && request.Limit > 0 {
		offset := request.Skip
		query = query.Offset(offset).Limit(request.Limit)
	}

	if err := query.Find(&campaigns).Error; err != nil {
		log.Error("Failed to fetch campaigns", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to fetch campaigns: %w", err)
	}

	go func() {
		cacheResponse := struct {
			Campaigns []models.Campaign `json:"campaigns"`
			Total     int64             `json:"total"`
		}{
			Campaigns: campaigns,
			Total:     total,
		}

		cacheJSON, err := json.Marshal(cacheResponse)
		if err != nil {
			log.Error("Failed to marshal cache response", zap.Error(err))
		} else {
			if err := redis.Set(ctx, cacheKey, string(cacheJSON), 15*time.Minute).Err(); err != nil {
				log.Error("Failed to cache query result", zap.Error(err))
			}
		}

		for _, campaign := range campaigns {
			campaignKey := fmt.Sprintf("campaign:Id:%s", campaign.ID.String())
			campaignJSON, err := json.Marshal(campaign)
			if err != nil {
				log.Error("Failed to marshal individual campaign for caching", zap.Error(err))
				continue
			}

			if err := redis.Set(ctx, campaignKey, campaignJSON, 30*time.Minute).Err(); err != nil {
				log.Error("Failed to set individual campaign in redis", zap.Error(err))
			}
		}

	}()

	return campaigns, total, nil
}
