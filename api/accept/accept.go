package accept

import (
	"campaign-service/business/accept"
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

func AcceptCampaign(ctx *gin.Context) {
	//get the lang
	lang, _ := ctx.Get(constants.LanguageString)

	//get the logger
	log := logger.GetLogger(ctx)

	var request models.AcceptCampaignRequest
	if validationErr := helperfunctions.ValidateRequestData(ctx, &request, binding.JSON); validationErr != nil {
		return
	}
	err := accept.AcceptCampaign(ctx, &request)

	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	//sent the success message
	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, "Campaign created successfully")
}
