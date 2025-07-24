package activity

import (
	"campaign-service/library/postgres"
	"campaign-service/logger"
	"campaign-service/models"
	"encoding/json"

	"go.uber.org/zap"
)

func ConsumeJoinActivity(participant map[string]interface{}) error {
	// get the logger
	log := logger.GetLoggerWithoutContext()

	// marshal the participant to json
	participantJSON, err := json.Marshal(participant)
	if err != nil {
		log.With(zap.Error(err)).Error("failed to marshal participant")
		return err
	}

	log.Info("Campaign activity join event received", zap.String("participant", string(participantJSON)))

	db := postgres.DB

	// start a transaction
	tx := db.Begin()

	// insert the participant into the db
	err = db.Model(&models.Participant{}).Create(&participant).Error
	if err != nil {
		log.With(zap.Error(err)).Error("failed to insert participant into db")
		tx.Rollback()
		return err
	}

	// commit the transaction
	tx.Commit()

	log.Info("Campaign activity join event consumed")
	return nil
}

func ConsumeLeaveActivity(participant map[string]interface{}) error {
	// get the logger
	log := logger.GetLoggerWithoutContext()

	// marshal the participant to json
	participantJSON, err := json.Marshal(participant)
	if err != nil {
		log.With(zap.Error(err)).Error("failed to marshal participant")
		return err
	}

	log.Info("Campaign activity join event received", zap.String("participant", string(participantJSON)))

	db := postgres.DB

	// start a transaction
	tx := db.Begin()

	// delete the participant from the db
	err = db.Model(&models.Participant{}).Where("id = ?", participant["id"]).Delete(&models.Participant{}).Error
	if err != nil {
		log.With(zap.Error(err)).Error("failed to delete participant from db")
		tx.Rollback()
		return err
	}

	// commit the transaction
	tx.Commit()

	log.Info("Campaign activity leave event consumed")
	return nil
}
