package nearby

import (
	"campaign-service/constants"
	"campaign-service/library/postgres"
	"campaign-service/logger"
	"campaign-service/models"
	"campaign-service/utils/helperfunctions"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetNearbyCampaigns(ctx *gin.Context, request *models.GetCampaignRequest) ([]*models.Campaign, int64, error) {
	log := logger.GetLoggerWithoutContext()

	userID := request.UserID
	if userID == "" {
		log.With(zap.Error(errors.New(constants.UserNotFoundMessage))).Error(constants.UserNotFoundMessage)
		return nil, 0, errors.New(constants.UserNotFoundMessage)
	}

	_, err := helperfunctions.ValidateUserExists(ctx, userID)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.UserNotFoundMessage)
		return nil, 0, err
	}

	latitudeStr := request.Latitude
	longitudeStr := request.Longitude
	radiusStr := request.Radius

	if latitudeStr == "" || longitudeStr == "" {
		return nil, 0, fmt.Errorf("latitude and longitude are required")
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid latitude format: %w", err)
	}

	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid longitude format: %w", err)
	}

	radius := 10.0
	if radiusStr != "" {
		if parsedRadius, err := strconv.ParseFloat(radiusStr, 64); err == nil && parsedRadius > 0 {
			radius = parsedRadius
		}
	}

	if latitude < -90 || latitude > 90 {
		return nil, 0, fmt.Errorf("latitude must be between -90 and 90")
	}
	if longitude < -180 || longitude > 180 {
		return nil, 0, fmt.Errorf("longitude must be between -180 and 180")
	}

	log.Info("Searching for nearby campaigns",
		zap.Float64("latitude", latitude),
		zap.Float64("longitude", longitude),
		zap.Float64("radius", radius))

	campaignIDs, err := helperfunctions.GetCampaignsInRadius(ctx, latitude, longitude, radius, "km")
	if err != nil {
		log.Error("Failed to get campaigns from spatial index", zap.Error(err))
		return nil, 0, fmt.Errorf("spatial search failed: %w", err)
	}

	if len(campaignIDs) == 0 {
		log.Info("No campaigns found in spatial index")
		return []*models.Campaign{}, 0, nil
	}

	log.Info("Found campaigns in spatial index", zap.Int("count", len(campaignIDs)))

	var campaigns []*models.Campaign
	db := postgres.DB

	query := db.Model(&models.Campaign{}).
		Preload("Location").
		Preload("StatusLogs").
		Where("id IN ?", campaignIDs)

	if request.Status != "" {
		query = query.Where("status = ?", request.Status)
	}

	if request.UserID != "" {
		query = query.Where("user_id = ?", request.UserID)
	}

	if request.MinPrice != 0 {
		query = query.Where("price >= ?", request.MinPrice)
	}

	if request.MaxPrice != 0 {
		query = query.Where("price <= ?", request.MaxPrice)
	}

	if request.City != "" {
		query = query.Joins("JOIN locations ON campaigns.id = locations.campaign_id").
			Where("locations.city = ?", request.City)
	}

	if request.State != "" {
		query = query.Joins("JOIN locations ON campaigns.id = locations.campaign_id").
			Where("locations.state = ?", request.State)
	}

	if request.Country != "" {
		query = query.Joins("JOIN locations ON campaigns.id = locations.campaign_id").
			Where("locations.country = ?", request.Country)
	}

	if request.StartDate != "" {
		query = query.Where("start_date >= ?", request.StartDate)
	}

	if request.EndDate != "" {
		query = query.Where("end_date <= ?", request.EndDate)
	}

	if request.SortBy != "" {
		sortOrder := "ASC"
		if request.SortOrder == "desc" {
			sortOrder = "DESC"
		}

		switch request.SortBy {
		case "start_date":
			query = query.Order(fmt.Sprintf("start_date %s", sortOrder))
		case "created_at":
			query = query.Order(fmt.Sprintf("created_at %s", sortOrder))
		case "price":
			query = query.Order(fmt.Sprintf("price %s", sortOrder))
		case "participants":
			query = query.Order(fmt.Sprintf("current_count %s", sortOrder))
		default:
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	if request.Limit > 0 {
		query = query.Limit(request.Limit)
	} else {
		query = query.Limit(20) // Default limit
	}

	if request.Skip > 0 {
		query = query.Offset(request.Skip)
	}

	if err := query.Find(&campaigns).Error; err != nil {
		log.Error("Failed to fetch campaigns from database", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to fetch campaigns: %w", err)
	}

	var total int64
	countQuery := db.Model(&models.Campaign{}).Where("id IN ?", campaignIDs)

	if request.Status != "" {
		countQuery = countQuery.Where("status = ?", request.Status)
	}
	if request.UserID != "" {
		countQuery = countQuery.Where("user_id = ?", request.UserID)
	}
	if request.MinPrice != 0 {
		countQuery = countQuery.Where("price >= ?", request.MinPrice)
	}
	if request.MaxPrice != 0 {
		countQuery = countQuery.Where("price <= ?", request.MaxPrice)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		log.Error("Failed to get total count", zap.Error(err))
		total = int64(len(campaigns))
	}

	log.Info("Nearby campaigns search completed",
		zap.Int("results", len(campaigns)),
		zap.Int64("total", total),
		zap.Float64("radius", radius))

	return campaigns, total, nil
}

func GetNearbyCampaignsSimple(ctx *gin.Context, latitude, longitude, radius float64, limit int) ([]*models.Campaign, error) {
	log := logger.GetLoggerWithoutContext()

	campaignIDs, err := helperfunctions.GetCampaignsInRadius(ctx, latitude, longitude, radius, "km")
	if err != nil {
		log.Error("Failed to get campaigns from spatial index", zap.Error(err))
		return nil, fmt.Errorf("spatial search failed: %w", err)
	}

	if len(campaignIDs) == 0 {
		return []*models.Campaign{}, nil
	}

	var campaigns []*models.Campaign
	db := postgres.DB

	query := db.Model(&models.Campaign{}).
		Preload("Location").
		Where("id IN ?", campaignIDs).
		Where("status = ?", "active").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(20)
	}

	if err := query.Find(&campaigns).Error; err != nil {
		log.Error("Failed to fetch campaigns from database", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch campaigns: %w", err)
	}

	log.Info("Simple nearby search completed", zap.Int("results", len(campaigns)))
	return campaigns, nil
}
