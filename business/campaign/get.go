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
	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the client name from the request
	userID := request.UserID
	if userID == "" {
		log.With(zap.Error(errors.New(constants.UserNotFoundMessage))).Error(constants.UserNotFoundMessage)
		return nil, 0, errors.New(constants.UserNotFoundMessage)
	}

	redis := redis_provider.Client

	// create redis cache key
	cacheKey := fmt.Sprintf(":user:%s", userID)
	indexKey := fmt.Sprintf("campaign:user:%s:index", userID)

	if request.ID != "" {
		// Try to get the campaign by ID directly from Redis
		campaignKey := fmt.Sprintf("campaign:user:%s:%s", userID, request.ID)
		campaignData, err := redis.Get(ctx, campaignKey).Result()
		if err == nil && campaignData != "" {
			var campaign models.Campaign
			if err := json.Unmarshal([]byte(campaignData), &campaign); err == nil {
				return []models.Campaign{campaign}, 1, nil
			}
		}
		// If not found in Redis, continue to DB logic below
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

	// Build the index key for the user
	indexKey = fmt.Sprintf("campaign:user:%s:index", userID)
	campaignIDs, err := redis.SMembers(ctx, indexKey).Result()
	if err == nil && len(campaignIDs) > 0 {
		var campaignKeys []string
		for _, id := range campaignIDs {
			campaignKeys = append(campaignKeys, fmt.Sprintf("campaign:user:%s:%s", userID, id))
		}
		redisResults, err := redis.MGet(ctx, campaignKeys...).Result()
		if err == nil && len(redisResults) > 0 {
			var campaigns []models.Campaign
			for _, result := range redisResults {
				if result == nil {
					continue
				}
				var campaign models.Campaign
				if err := json.Unmarshal([]byte(result.(string)), &campaign); err == nil {
					campaigns = append(campaigns, campaign)
				}
			}
			if len(campaigns) > 0 {
				log.Info("***************get the campaigns form redis****************")
				return campaigns, int64(len(campaigns)), nil
			}
		}
	}

	// 	db operation starts
	db := postgres.DB
	query := db.Model(&models.Campaign{}).Preload("Location").Preload("StatusLogs").Where("user_id = ?", userID)
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

	// Handle tags filter - search campaigns that contain any of the specified tags
	if len(request.Tags) > 0 {
		query = query.Where("tags ?| ?", request.Tags)
	}

	// Add sorting
	if request.SortBy != "" {
		sortOrder := "ASC"
		if request.SortOrder != "" {
			sortOrder = strings.ToUpper(request.SortOrder)
		}
		query = query.Order(fmt.Sprintf("%s %s", request.SortBy, sortOrder))
	}

	// Execute the query
	var campaigns []models.Campaign
	var total int64

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		log.Error("Failed to count campaigns", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count campaigns: %w", err)
	}

	// Get paginated results
	if request.Skip > 0 && request.Limit > 0 {
		offset := request.Skip
		query = query.Offset(offset).Limit(request.Limit)
	}

	if err := query.Find(&campaigns).Error; err != nil {
		log.Error("Failed to fetch campaigns", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to fetch campaigns: %w", err)
	}
	log.Info("***************get the campaigns form postgres****************")

	go func() {
		// cache the result on redis
		campaignsJSON, err := json.Marshal(campaigns)
		if err != nil {
			log.Error("Failed to marshal campaigns for caching", zap.Error(err))
		} else {
			redis.Set(ctx, cacheKey, string(campaignsJSON), 15*time.Minute)
		}

		// For each campaign, cache individually and add to user's index set
		for _, campaign := range campaigns {
			campaignKey := fmt.Sprintf("campaign:user:%s:%s", userID, campaign.ID.String())
			campaignJSON, err := json.Marshal(campaign)
			if err != nil {
				log.Error("Failed to marshal individual campaign for caching", zap.Error(err))
				continue
			}
			if err := redis.Set(ctx, campaignKey, campaignJSON, 5*time.Minute).Err(); err != nil {
				log.Error("Failed to set individual campaign in redis", zap.Error(err))
			}
			if err := redis.SAdd(ctx, indexKey, campaign.ID.String()).Err(); err != nil {
				log.Error("Failed to add campaign ID to user index set", zap.Error(err))
			}
		}

		log.Info("***************set the campaigns to redis****************")
	}()

	return campaigns, total, nil
}
