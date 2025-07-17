package middleware

import (
	"campaign-service/constants"
	"campaign-service/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

var UserDatafromGataway models.UserDataFromAPIGateWay

func KecyalokMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Set data from authorization header
		authToken := ctx.GetHeader(constants.AuthorizationHeader)
		json.Unmarshal([]byte(authToken), &UserDatafromGataway)

		if UserDatafromGataway.AuthCompleted {
			ctx.AbortWithStatusJSON(http.StatusFailedDependency, constants.DirectAccessNotAllowed)
		}
	}
}
