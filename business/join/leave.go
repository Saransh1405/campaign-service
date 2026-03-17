package join

import (
	"campaign-service/constants"
	"campaign-service/library/kafka/activity"
	"campaign-service/library/postgres"
	"campaign-service/library/redis_provider"
	"campaign-service/logger"
	"campaign-service/models"
	"campaign-service/utils/helperfunctions"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func LeaveCampaign(ctx context.Context, request *models.LeaveCampaignRequest) error {
	log := logger.GetLoggerWithoutContext()

	userID := request.UserID
	if userID == "" {
		log.With(zap.Error(errors.New(constants.UserNotFoundMessage))).Error(constants.UserNotFoundMessage)
		return errors.New(constants.UserNotFoundMessage)
	}

	user, err := helperfunctions.ValidateUserExists(ctx, userID)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.UserNotFoundMessage)
		return err
	}

	if !user.EmailVerified {
		log.With(zap.Error(errors.New(constants.UserNotVerifiedMessage))).Error(constants.UserNotVerifiedMessage)
		return errors.New(constants.UserNotVerifiedMessage)
	}

	db := postgres.DB
	var campaign models.Campaign
	campaignKey := fmt.Sprintf("campaign:user:%s:%s", userID, request.CampaignID)
	campaignData, err := redis_provider.Client.Get(ctx, campaignKey).Result()
	if err == redis.Nil {
		err := db.Model(&models.Campaign{}).Where("id = ?", request.CampaignID).First(&campaign).Error
		if err != nil {
			log.With(zap.Error(err)).Error(constants.CampaignNotFoundMessage)
			return errors.New(constants.CampaignNotFoundMessage)
		}
	} else if err != nil {
		log.With(zap.Error(err)).Error("Redis error")
		return err
	} else {
		if err := json.Unmarshal([]byte(campaignData), &campaign); err != nil {
			log.With(zap.Error(err)).Error("Failed to unmarshal campaign from cache")
			return err
		}
	}

	if campaign.CurrentCount >= campaign.MaxParticipants {
		log.With(zap.Error(errors.New(constants.CampaignFullMessage))).Error(constants.CampaignFullMessage)
		return errors.New(constants.CampaignFullMessage)
	}

	err = db.Model(&models.Participant{}).Where("user_id = ? AND campaign_id = ?", userID, request.CampaignID).First(&models.Participant{}).Error
	if err == gorm.ErrRecordNotFound {
		log.With(zap.Error(err)).Error("failed to check if user is already in the campaign")
		return errors.New(constants.UserNotParticipantMessage)
	}

	go func() {
		var participant models.Participant
		err = db.Model(&models.Participant{}).Where("user_id = ? AND campaign_id = ?", userID, request.CampaignID).First(&participant).Error
		if err != nil {
			log.With(zap.Error(err)).Error("failed to read participant from db")
			return
		}
		err = SendCampaignLeaveActivityToKafka(participant, string(models.CampaignActivity))
		if err != nil {
			log.With(zap.Error(err)).Error("failed to send campaign activity to kafka")
		}
	}()

	go func() {
		statusLog := models.StatusLogs{
			ID:             uuid.New(),
			CampaignID:     campaign.ID,
			Status:         models.Active,
			ActionByUserId: userID,
			Notes:          "user left the campaign",
			Timestamp:      time.Now().UnixMilli(),
		}

		err = db.Model(&models.StatusLogs{}).Create(&statusLog).Error
		if err != nil {
			log.With(zap.Error(err)).Error("failed to insert status logs into db")
		}

		log.Info("Campaign activity status logs inserted into db")
	}()

	go func() {
		if err := db.WithContext(ctx).
			Model(&models.Campaign{}).
			Where("id = ?", request.CampaignID).
			Update("current_count", gorm.Expr("current_count - 1")).Error; err != nil {
			log.With(zap.Error(err)).Error("Failed to decrement campaign current_count")
		}
		log.Info("Campaign count decreased")
	}()

	go func() {
		backgroundContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := helperfunctions.InvalidateAllCampaignUserCache(backgroundContext, userID); err != nil {
			log.With(zap.Error(err)).Error("Failed to invalidate campaign user cache")
		}
	}()

	log.Info("Campaign activity leave event published to kafka")
	return nil
}

func SendCampaignLeaveActivityToKafka(participant models.Participant, topic string) error {
	log := logger.GetLoggerWithoutContext()

	participantData := map[string]interface{}{
		"id":          participant.ID,
		"user_id":     participant.UserID,
		"campaign_id": participant.CampaignID,
		"joined_at":   participant.JoinedAt,
		"status":      participant.Status,
		"payment_id":  participant.PaymentID,
		"created_at":  time.Now().UnixMilli(),
		"updated_at":  time.Now().UnixMilli(),
	}

	payload := &models.ActivityEvent{
		Participant:      participantData,
		Action:           "leave",
		EventPublishTime: time.Now().UnixMilli(),
	}

	err := activity.SendActivityDataToKafka(payload, string(models.CampaignActivity))
	if err != nil {
		return err
	}

	log.Info("Campaign activity join event published to kafka")
	return nil
}
