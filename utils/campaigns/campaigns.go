package campaigns

import (
	"campaign-service/library/postgres"
	"campaign-service/logger"
	"campaign-service/models"
	"encoding/json"
	"time"

	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func InsertIntoPostgres(campaignData models.Campaign) error {
	log := logger.GetLoggerWithoutContext()

	db := postgres.DB

	// start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		log.Error("Error in starting transaction", zap.Error(tx.Error))
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// create the campaign object
	campaign := models.Campaign{
		ID:              campaignData.ID,
		UserID:          campaignData.UserID,
		Name:            campaignData.Name,
		Description:     campaignData.Description,
		Type:            campaignData.Type,
		ImageURL:        campaignData.ImageURL,
		DisplayName:     campaignData.DisplayName,
		Price:           campaignData.Price,
		MinParticipants: campaignData.MinParticipants,
		MaxParticipants: campaignData.MaxParticipants,
		StartDate:       campaignData.StartDate,
		EndDate:         campaignData.EndDate,
		Currency:        campaignData.Currency,
		Tags:            json.RawMessage(campaignData.Tags),
		Status:          campaignData.Status,
		IsPublic:        campaignData.IsPublic,
		CreatedAt:       campaignData.CreatedAt,
		UpdatedAt:       campaignData.UpdatedAt,
	}
	if result := db.Save(&campaign); result.Error != nil {
		tx.Rollback()
		log.Error("Error in saving campaign", zap.Error(result.Error))
		return fmt.Errorf("failed to save campaign: %w", result.Error)
	}

	// create the location object
	location := models.Location{
		ID:         uuid.New(),
		CampaignID: campaignData.ID,
		Name:       campaignData.Location.Name,
		Address:    campaignData.Location.Address,
		Latitude:   campaignData.Location.Latitude,
		Longitude:  campaignData.Location.Longitude,
		City:       campaignData.Location.City,
		State:      campaignData.Location.State,
		Country:    campaignData.Location.Country,
		ZipCode:    campaignData.Location.ZipCode,
		CreatedAt:  campaignData.Location.CreatedAt,
		UpdatedAt:  campaignData.Location.UpdatedAt,
	}
	if result := db.Save(&location); result.Error != nil {
		tx.Rollback()
		log.Error("Error in saving location", zap.Error(result.Error))
		return fmt.Errorf("failed to save location: %w", result.Error)
	}

	// create the status logs object
	statusLogs := models.StatusLogs{
		ID:             uuid.New(),
		CampaignID:     campaignData.ID,
		Status:         models.Active, // Use enum instead of string
		Notes:          "Campaign created",
		ActionByUserId: campaignData.UserID,
		Timestamp:      time.Now().UnixMilli(), // Use int64 timestamp
	}
	if result := db.Create(&statusLogs); result.Error != nil {
		tx.Rollback()
		log.Error("Error in creating status logs", zap.Error(result.Error))
		return fmt.Errorf("failed to create status logs: %w", result.Error)
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("Error in committing transaction", zap.Error(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info("Campaign event published to kafka")

	return nil
}

func UpdateCampaignInPostgres(campaignData map[string]interface{}, updateFields map[string]interface{}) error {
	log := logger.GetLoggerWithoutContext()

	// Initialize updateFields if it's nil
	if updateFields == nil {
		updateFields = make(map[string]interface{})
		fmt.Printf("updateFields: %v\n", updateFields)
	}

	jsonBytes, err := json.Marshal(campaignData)
	if err != nil {
		log.Error("Error in marshalling campaign data", zap.Error(err))
		return fmt.Errorf("failed to marshal campaign data: %w", err)
	}

	var campaign models.Campaign
	if err := json.Unmarshal(jsonBytes, &campaign); err != nil {
		log.Error("Error in unmarshalling campaign data", zap.Error(err))
		return fmt.Errorf("failed to unmarshal campaign data: %w", err)
	}

	db := postgres.DB

	// start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		log.Error("Error in starting transaction", zap.Error(tx.Error))
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// Remove location from updateFields since it's handled separately
	delete(updateFields, "location")

	// update the tags as jsonb
	updateFields["tags"] = json.RawMessage(campaign.Tags)

	// Add updated_at timestamp
	updateFields["updated_at"] = time.Now().UnixMilli()

	// Update the campaign with only the changed fields
	if result := db.Model(&models.Campaign{}).Where("id = ?", campaign.ID).Updates(updateFields); result.Error != nil {
		log.Error("Error in updating campaign", zap.Error(result.Error))
		tx.Rollback()
		return fmt.Errorf("failed to update campaign: %w", result.Error)
	}

	// Update location if provided
	if campaign.Location != nil {
		locationUpdate := map[string]interface{}{
			"name":       campaign.Location.Name,
			"address":    campaign.Location.Address,
			"latitude":   campaign.Location.Latitude,
			"longitude":  campaign.Location.Longitude,
			"city":       campaign.Location.City,
			"state":      campaign.Location.State,
			"country":    campaign.Location.Country,
			"zip_code":   campaign.Location.ZipCode,
			"updated_at": time.Now().UnixMilli(),
		}

		if result := db.Model(&models.Location{}).Where("campaign_id = ?", campaign.ID).Updates(locationUpdate); result.Error != nil {
			log.Error("Error in updating location", zap.Error(result.Error))
			tx.Rollback()
			return fmt.Errorf("failed to update location: %w", result.Error)
		}
	}

	// Add status log entry
	statusLogs := models.StatusLogs{
		ID:             uuid.New(),
		CampaignID:     campaign.ID,
		Status:         models.Active, // Use lowercase string to match enum
		Notes:          "Campaign updated",
		ActionByUserId: campaign.UserID,
		Timestamp:      time.Now().UnixMilli(),
	}

	if result := db.Create(&statusLogs); result.Error != nil {
		log.Error("Error in creating status logs", zap.Error(result.Error))
		tx.Rollback()
		return fmt.Errorf("failed to create status logs: %w", result.Error)
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("Error in committing transaction", zap.Error(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info("Campaign updated successfully in database")
	return nil
}
