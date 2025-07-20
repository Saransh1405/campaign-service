package redis_provider

import (
	"campaign-service/logger"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// CampaignSpatialIndex is the single spatial index for all campaigns
const CampaignSpatialIndex = "campaigns:spatial"

// AddCampaignToSpatialIndex adds a campaign to the spatial index
func AddCampaignToSpatialIndex(ctx context.Context, campaignID string, latitude, longitude float64) error {
	log := logger.GetLoggerWithoutContext()

	// Create Redis GEO location
	geoLocation := &redis.GeoLocation{
		Name:      campaignID,
		Longitude: longitude,
		Latitude:  latitude,
	}

	// Add to spatial index using Redis client directly
	if err := Client.GeoAdd(ctx, CampaignSpatialIndex, geoLocation).Err(); err != nil {
		return fmt.Errorf("failed to add campaign to spatial index: %w", err)
	}

	log.Info("campaign added to spatial index",
		zap.String("campaignID", campaignID),
		zap.Float64("latitude", latitude),
		zap.Float64("longitude", longitude))

	return nil
}

// RemoveCampaignFromSpatialIndex removes a campaign from the spatial index
func RemoveCampaignFromSpatialIndex(ctx context.Context, campaignID string) error {
	log := logger.GetLoggerWithoutContext()

	if err := Client.ZRem(ctx, CampaignSpatialIndex, campaignID).Err(); err != nil {
		return fmt.Errorf("failed to remove campaign from spatial index: %w", err)
	}

	log.Info("campaign removed from spatial index", zap.String("campaignID", campaignID))
	return nil
}

// GetCampaignsInRadius gets campaign IDs within a specified radius
func GetCampaignsInRadius(ctx context.Context, latitude, longitude, radius float64, unit string) ([]string, error) {
	log := logger.GetLoggerWithoutContext()

	if unit == "" {
		unit = "km"
	}

	query := &redis.GeoRadiusQuery{
		Radius:      radius,
		Unit:        unit,
		WithCoord:   false,
		WithDist:    false,
		WithGeoHash: false,
		Count:       100,   // Limit results
		Sort:        "ASC", // Sort by distance
	}

	// Perform spatial search
	locations, err := Client.GeoRadius(ctx, CampaignSpatialIndex, longitude, latitude, query).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to search campaigns by radius: %w", err)
	}

	// Extract campaign IDs
	var campaignIDs []string
	for _, location := range locations {
		campaignIDs = append(campaignIDs, location.Name)
	}

	log.Info("spatial search completed",
		zap.Int("results", len(campaignIDs)),
		zap.Float64("radius", radius),
		zap.String("unit", unit))

	return campaignIDs, nil
}

// GetSpatialIndexStats returns the number of campaigns in the spatial index
func GetSpatialIndexStats(ctx context.Context) (int64, error) {
	count, err := Client.ZCard(ctx, CampaignSpatialIndex).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get spatial index stats: %w", err)
	}

	return count, nil
}

// ClearSpatialIndex clears the spatial index
func ClearSpatialIndex(ctx context.Context) error {
	log := logger.GetLoggerWithoutContext()

	if err := Client.Del(ctx, CampaignSpatialIndex).Err(); err != nil {
		return fmt.Errorf("failed to clear spatial index: %w", err)
	}

	log.Info("spatial index cleared")
	return nil
}
