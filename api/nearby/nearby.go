package nearby

import (
	"campaign-service/business/nearby"
	"campaign-service/constants"
	"campaign-service/logger"
	"campaign-service/models"
	"campaign-service/utils"
	"campaign-service/utils/helperfunctions"
	"campaign-service/utils/localization"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetCampaign(ctx *gin.Context) {
	//get the lang
	lang, _ := ctx.Get(constants.LanguageString)

	//get the logger
	log := logger.GetLogger(ctx)

	var request models.GetCampaignRequest
	if validationErr := helperfunctions.ValidateRequestDataParam(ctx, &request); validationErr != nil {
		return
	}

	result, count, err := nearby.GetNearbyCampaigns(ctx, &request)
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

	//sent the success message
	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusWithData(ctx, constants.IsString, successMessage, result, count)
}
