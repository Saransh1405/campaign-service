package campaign

import (
	"campaign-service/constants"
	"campaign-service/library/kafka/campaign"
	"campaign-service/library/postgres"
	"campaign-service/library/redis_provider"
	"campaign-service/logger"
	"campaign-service/models"
	"campaign-service/utils/helperfunctions"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func UpdateCampaign(ctx *gin.Context, campaign *models.UpdateCampaignRequest) error {
	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the client name from the request
	userID := campaign.UserID
	if userID == "" {
		log.With(zap.Error(errors.New(constants.UserNotFoundMessage))).Error(constants.UserNotFoundMessage)
		return errors.New(constants.UserNotFoundMessage)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 5)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Time validation
		currentTime := time.Now().UnixMilli()

		if campaign.StartDate != nil && campaign.EndDate != nil {

			// Validate start date is not in the past
			if *campaign.StartDate < currentTime {
				log.With(zap.Error(errors.New(constants.InvalidStartDateMessage))).Error(constants.InvalidStartDateMessage)
				errCh <- errors.New(constants.InvalidStartDateMessage)
			}

			// Validate end date is not in the past
			if *campaign.EndDate < currentTime {
				log.With(zap.Error(errors.New(constants.InvalidEndDateMessage))).Error(constants.InvalidEndDateMessage)
				errCh <- errors.New(constants.InvalidEndDateMessage)
			}

			// Validate end date is after start date
			if *campaign.EndDate <= *campaign.StartDate {
				log.With(zap.Error(errors.New(constants.EndDateBeforeStartDateMessage))).Error(constants.EndDateBeforeStartDateMessage)
				errCh <- errors.New(constants.EndDateBeforeStartDateMessage)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if campaign.MaxParticipants != nil && campaign.MinParticipants != nil {
			// Validate max participants is greater than min participants
			if *campaign.MaxParticipants < *campaign.MinParticipants {
				log.With(zap.Error(errors.New(constants.MaxParticipantsLessThanMinParticipants))).Error(constants.MaxParticipantsLessThanMinParticipants)
				errCh <- errors.New(constants.MaxParticipantsLessThanMinParticipants)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if campaign.Price != nil {
			// Validate price is greater than 0
			if *campaign.Price <= 0 {
				log.With(zap.Error(errors.New(constants.PriceMustBeGreaterThanZero))).Error(constants.PriceMustBeGreaterThanZero)
				errCh <- errors.New(constants.PriceMustBeGreaterThanZero)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Validate name is valid check in db if it already exists
		if *campaign.Name != "" {
			db := postgres.DB
			campaignModel := db.Model(&models.Campaign{}).Where("name = ? AND id != ?", *campaign.Name, *campaign.ID).First(&models.Campaign{})
			if campaignModel.Error == nil {
				log.With(zap.Error(errors.New(constants.CampaignNameAlreadyExistsMessage))).Error(constants.CampaignNameAlreadyExistsMessage)
				errCh <- errors.New(constants.CampaignNameAlreadyExistsMessage)
			}
		}
	}()

	// Wait for all validations to finish
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Return the first error found
	for err := range errCh {
		log.With(zap.Error(err)).Error("Validation failed")
		return err
	}

	// channel to store the campaign model
	campaignModelCh := make(chan *models.Campaign, 1)
	camapignGetErrCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		// get the campaign from the database
		if campaign.ID != nil {
			var existingCampaign models.Campaign

			redis := redis_provider.Client
			campaignKey := fmt.Sprintf("campaign:user:%s:%s", userID, *campaign.ID)
			campaignData, err := redis.Get(ctx, campaignKey).Result()
			if err == nil && campaignData != "" {
				// Unmarshal the JSON into the struct
				if err := json.Unmarshal([]byte(campaignData), &existingCampaign); err != nil {
					log.With(zap.Error(err)).Error("Failed to unmarshal campaign from redis")
				}
				log.Info("***************get the campaign from redis****************")
			} else {
				// Not found in Redis, get from DB
				db := postgres.DB
				campaignModel := db.Model(&models.Campaign{}).Where("id = ?", *campaign.ID).First(&existingCampaign)
				if campaignModel.Error != nil {
					log.With(zap.Error(errors.New(constants.CampaignNotFoundMessage))).Error(constants.CampaignNotFoundMessage)
					camapignGetErrCh <- errors.New(constants.CampaignNotFoundMessage)
				}
				log.Info("***************get the campaign from postgres****************")
			}

			campaignModelCh <- &existingCampaign
		}
	}()

	go func() {
		wg.Wait()
		close(campaignModelCh)
		close(camapignGetErrCh)
	}()

	for err := range camapignGetErrCh {
		log.With(zap.Error(err)).Error("Error in getting campaign")
		return err
	}

	// get the campaign model
	newCampaign := <-campaignModelCh

	updateFields := make(map[string]interface{})

	if campaign.Name != nil {
		newCampaign.Name = *campaign.Name
		updateFields["name"] = *campaign.Name
	}

	if campaign.Description != nil {
		newCampaign.Description = *campaign.Description
		updateFields["description"] = *campaign.Description
	}

	if campaign.ImageURL != nil {
		newCampaign.ImageURL = *campaign.ImageURL
		updateFields["image_url"] = *campaign.ImageURL
	}

	if campaign.DisplayName != nil {
		newCampaign.DisplayName = *campaign.DisplayName
		updateFields["display_name"] = *campaign.DisplayName
	}

	if campaign.StartDate != nil {
		newCampaign.StartDate = *campaign.StartDate
		updateFields["start_date"] = *campaign.StartDate
	}

	if campaign.EndDate != nil {
		newCampaign.EndDate = *campaign.EndDate
		updateFields["end_date"] = *campaign.EndDate
	}

	if campaign.MaxParticipants != nil {
		newCampaign.MaxParticipants = *campaign.MaxParticipants
		updateFields["max_participants"] = *campaign.MaxParticipants
	}

	if campaign.MinParticipants != nil {
		newCampaign.MinParticipants = *campaign.MinParticipants
		updateFields["min_participants"] = *campaign.MinParticipants
	}

	if campaign.Price != nil {
		newCampaign.Price = *campaign.Price
		updateFields["price"] = *campaign.Price
	}

	if campaign.IsPublic != nil {
		newCampaign.IsPublic = *campaign.IsPublic
		updateFields["is_public"] = *campaign.IsPublic
	}

	if campaign.Tags != nil {
		newCampaign.Tags = json.RawMessage(campaign.Tags)
		updateFields["tags"] = json.RawMessage(campaign.Tags)
	}

	if campaign.Location != nil {
		newCampaign.Location = campaign.Location
		updateFields["location"] = campaign.Location
	}

	if campaign.Status != nil {
		newCampaign.Status = *campaign.Status
		updateFields["status"] = *campaign.Status
	}

	// Validate campaign.ID before parsing to avoid panic
	if campaign.ID == nil || *campaign.ID == "" {
		log.Error("campaign.ID is nil or empty, cannot parse UUID")
		return errors.New("invalid campaign ID: cannot be nil or empty")
	}
	parsedID, err := uuid.Parse(*campaign.ID)
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to parse campaign ID as UUID")
		return fmt.Errorf("invalid campaign ID: %w", err)
	}
	newCampaign.ID = parsedID

	go func() {
		if err := PublishUpdateDataToKafka(*newCampaign, updateFields, string(models.UpdateCampaignEvent)); err != nil {
			log.With(zap.Error(err)).Error("Failed to publish campaign event to kafka")
		}
	}()

	go func() {
		backgroundContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := helperfunctions.InvalidateAllCampaignUserCache(backgroundContext, userID); err != nil {
			log.With(zap.Error(err)).Error("Failed to invalidate campaign user cache")
		}
	}()

	// Re-cache the updated campaign and update the user's index set
	go func() {
		backgroundContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		campaignKey := fmt.Sprintf("campaign:user:%s:%s", userID, *campaign.ID)
		campaignJSON, err := json.Marshal(newCampaign)
		if err != nil {
			log.With(zap.Error(err)).Error("Failed to marshal updated campaign for Redis")
			return
		}
		if err := redis_provider.Client.Set(backgroundContext, campaignKey, campaignJSON, 15*time.Minute).Err(); err != nil {
			log.With(zap.Error(err)).Error("Failed to set updated campaign in Redis")
		}
		if err := helperfunctions.AddCampaignToUserIndex(backgroundContext, userID, *campaign.ID); err != nil {
			log.With(zap.Error(err)).Error("Failed to add campaign ID to user index set")
		}
	}()

	return nil
}

func PublishUpdateDataToKafka(campaignEvent models.Campaign, updateFields map[string]interface{}, topic string) error {
	log := logger.GetLoggerWithoutContext()

	campaignData := map[string]interface{}{
		"id":               campaignEvent.ID,
		"name":             campaignEvent.Name,
		"description":      campaignEvent.Description,
		"image_url":        campaignEvent.ImageURL,
		"display_name":     campaignEvent.DisplayName,
		"price":            campaignEvent.Price,
		"min_participants": campaignEvent.MinParticipants,
		"max_participants": campaignEvent.MaxParticipants,
		"start_date":       campaignEvent.StartDate,
		"end_date":         campaignEvent.EndDate,
		"currency":         campaignEvent.Currency,
		"tags":             campaignEvent.Tags,
		"location":         campaignEvent.Location,
		"type":             campaignEvent.Type,
		"status":           campaignEvent.Status,
		"created_at":       campaignEvent.CreatedAt,
		"updated_at":       campaignEvent.UpdatedAt,
		"user_id":          campaignEvent.UserID,
	}

	fmt.Printf("updateFields: %v\n", updateFields)

	payload := &models.CampaignEvent{
		Campaign:         campaignData,
		UpdateFields:     updateFields,
		EventType:        string(models.CreateCampaignEvent),
		EventPublishTime: time.Now().UnixMilli(),
	}

	err := campaign.SendDataToKafka(payload, string(models.UpdateCampaignEvent))
	if err != nil {
		return err
	}

	log.Info("Campaign update event published to kafka")
	return nil
}
