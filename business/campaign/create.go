package campaign

import (
	"campaign-service/constants"
	"campaign-service/library/kafka/campaign"
	"campaign-service/library/postgres"
	"campaign-service/logger"
	"campaign-service/models"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateCampaign(ctx *gin.Context, campaign *models.CreateCampaignRequest) error {
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

		// Validate start date is not in the past
		if campaign.StartDate < currentTime {
			log.With(zap.Error(errors.New(constants.InvalidStartDateMessage))).Error(constants.InvalidStartDateMessage)
			errCh <- errors.New(constants.InvalidStartDateMessage)
		}

		// Validate end date is not in the past
		if campaign.EndDate < currentTime {
			log.With(zap.Error(errors.New(constants.InvalidEndDateMessage))).Error(constants.InvalidEndDateMessage)
			errCh <- errors.New(constants.InvalidEndDateMessage)
		}

		// Validate end date is after start date
		if campaign.EndDate <= campaign.StartDate {
			log.With(zap.Error(errors.New(constants.EndDateBeforeStartDateMessage))).Error(constants.EndDateBeforeStartDateMessage)
			errCh <- errors.New(constants.EndDateBeforeStartDateMessage)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Validate max participants is greater than min participants
		if campaign.MaxParticipants < campaign.MinParticipants {
			log.With(zap.Error(errors.New(constants.MaxParticipantsLessThanMinParticipants))).Error(constants.MaxParticipantsLessThanMinParticipants)
			errCh <- errors.New(constants.MaxParticipantsLessThanMinParticipants)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Validate price is greater than 0
		if campaign.Price <= 0 {
			log.With(zap.Error(errors.New(constants.PriceMustBeGreaterThanZero))).Error(constants.PriceMustBeGreaterThanZero)
			errCh <- errors.New(constants.PriceMustBeGreaterThanZero)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Validate name is valid check in db if it already exists
		if campaign.Name != "" {
			db := postgres.DB
			campaignCol := db.Model(&models.Campaign{}).Where("name = ?", campaign.Name).First(&models.Campaign{})
			if campaignCol.Error == nil {
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

	//create the campaign
	createdAt := time.Now().UnixMilli()
	campaignModel := models.Campaign{
		ID:              uuid.New(),
		UserID:          userID,
		Name:            campaign.Name,
		Description:     campaign.Description,
		ImageURL:        campaign.ImageURL,
		DisplayName:     campaign.DisplayName,
		Price:           campaign.Price,
		MinParticipants: campaign.MinParticipants,
		MaxParticipants: campaign.MaxParticipants,
		StartDate:       campaign.StartDate,
		EndDate:         campaign.EndDate,
		Currency:        campaign.Currency,
		Tags:            campaign.Tags,
		Location:        &campaign.Location,
		Type:            campaign.Type,
		Status:          models.CampaignStatusActive,
		IsPublic:        campaign.IsPublic,
		CreatedAt:       createdAt,
		UpdatedAt:       createdAt,
	}

	go func() {
		if err := PublishDataToKafka(campaignModel, string(models.CreateCampaignEvent)); err != nil {
			log.With(zap.Error(err)).Error("Failed to publish campaign event to kafka")
		}
	}()

	return nil
}

func PublishDataToKafka(campaignEvent models.Campaign, topic string) error {
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

	payload := &models.CampaignEvent{
		Campaign:         campaignData,
		EventType:        string(models.CreateCampaignEvent),
		EventPublishTime: time.Now().UnixMilli(),
	}

	err := campaign.SendDataToKafka(payload, string(models.CreateCampaignEvent))
	if err != nil {
		return err
	}

	log.Info("Campaign insert event published to kafka")
	return nil
}
