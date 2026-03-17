package campaign

import (
	"campaign-service/business/campaign"
	"campaign-service/constants"
	"campaign-service/logger"
	"campaign-service/models"
	"campaign-service/utils"
	"campaign-service/utils/helperfunctions"
	"campaign-service/utils/localization"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

func CreateCampaign(ctx *gin.Context) {
	lang, _ := ctx.Get(constants.LanguageString)

	log := logger.GetLogger(ctx)

	var request models.CreateCampaignRequest
	if validationErr := helperfunctions.ValidateRequestData(ctx, &request, binding.JSON); validationErr != nil {
		return
	}
	err := campaign.CreateCampaign(ctx, &request)

	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, "Campaign created successfully")
}

func UpdateCampaign(ctx *gin.Context) {
	lang, _ := ctx.Get(constants.LanguageString)

	log := logger.GetLogger(ctx)

	var request models.UpdateCampaignRequest
	if validationErr := helperfunctions.ValidateRequestData(ctx, &request, binding.JSON); validationErr != nil {
		return
	}
	err := campaign.UpdateCampaign(ctx, &request)

	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, "Campaign created successfully")
}

func GetCampaign(ctx *gin.Context) {
	lang, _ := ctx.Get(constants.LanguageString)

	log := logger.GetLogger(ctx)

	var request models.GetCampaignRequest
	if validationErr := helperfunctions.ValidateRequestDataParam(ctx, &request); validationErr != nil {
		return
	}

	result, count, err := campaign.GetCampaign(ctx, &request)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	if result == nil || count == 0 {
		msg := localization.GetMessage(lang, constants.NoContentErrr, nil)
		utils.SendNoContentError(ctx, msg, constants.NoContentErrr, constants.IsString, err)
		return
	}

	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusWithData(ctx, constants.IsString, successMessage, result, count)
}
