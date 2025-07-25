package accept

import (
	"campaign-service/constants"
	"campaign-service/library/postgres"
	"campaign-service/logger"
	"campaign-service/models"
	"campaign-service/utils/helperfunctions"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func AcceptCampaign(ctx *gin.Context, request *models.AcceptCampaignRequest) error {
	// get the logger
	log := logger.GetLogger(ctx)

	//get the client name from the request
	userID := request.UserID
	if userID == "" {
		log.With(zap.Error(errors.New(constants.UserNotFoundMessage))).Error(constants.UserNotFoundMessage)
		return errors.New(constants.UserNotFoundMessage)
	}

	// validate the user exists
	user, err := helperfunctions.ValidateUserExists(ctx, userID)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.UserNotFoundMessage)
		return errors.New(constants.UserNotFoundMessage)
	}

	if !user.EmailVerified {
		log.With(zap.Error(errors.New(constants.UserNotVerifiedMessage))).Error(constants.UserNotVerifiedMessage)
		return errors.New(constants.UserNotVerifiedMessage)
	}

	// get the campaign
	db := postgres.DB

	// check with campaign exists
	campaignCol := db.Model(&models.Campaign{}).Where("id = ?", request.CampaignID).First(&models.Campaign{})
	if campaignCol.Error == gorm.ErrRecordNotFound {
		log.With(zap.Error(errors.New(constants.CampaignNotFoundMessage))).Error(constants.CampaignNotFoundMessage)
		return errors.New(constants.CampaignNotFoundMessage)
	}

	// check with user is already a participant
	participantCol := db.Model(&models.Participant{}).Where("user_id = ? AND campaign_id = ?", userID, request.CampaignID).First(&models.Participant{})
	if participantCol.Error == gorm.ErrRecordNotFound {
		log.With(zap.Error(errors.New(constants.UserNotParticipantMessage))).Error(constants.UserNotParticipantMessage)
		return errors.New(constants.UserNotParticipantMessage)
	}

	// update the participant status
	var status models.ParticipantStatus
	if request.Accept {
		status = models.ParticipantStatusActive
	} else {
		status = models.ParticipantStatusRejected
	}

	err = db.Model(&models.Participant{}).Where("user_id = ? AND campaign_id = ?", userID, request.CampaignID).Update("status", status).Error
	if err != nil {
		log.With(zap.Error(err)).Error(constants.FailedToUpdateParticipantStatusMessage)
		return errors.New(constants.FailedToUpdateParticipantStatusMessage)
	}

	go func() {
		// update the campaign current count
		if request.Accept {
			if err := db.WithContext(ctx).
				Model(&models.Campaign{}).
				Where("id = ?", request.CampaignID).
				Update("current_count", gorm.Expr("current_count + 1")).Error; err != nil {
				log.With(zap.Error(err)).Error("Failed to increment campaign current_count")
			}
			log.Info("Campaign count increased")
		}
	}()

	go func() {
		backgroundContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := helperfunctions.InvalidateAllCampaignUserCache(backgroundContext, userID); err != nil {
			log.With(zap.Error(err)).Error("Failed to invalidate campaign user cache")
		}
	}()

	return nil
}
