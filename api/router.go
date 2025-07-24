package api

import (
	"campaign-service/api/campaign"
	"campaign-service/api/join"
	"campaign-service/api/nearby"
	"campaign-service/constants"
	"campaign-service/utils"
	"campaign-service/utils/localization"
	"campaign-service/utils/middleware"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// GetRouter is used to get the router configured with the middlewares and the routes.
func GetRouter(localizationMiddleware gin.HandlerFunc, loggerMiddleware gin.HandlerFunc, applicationConfig *viper.Viper) *gin.Engine {
	router := gin.New()

	router.Use(localizationMiddleware)
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware)

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))

	middlewareFunc := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "DELETE", "PATCH", "PUT"},
		AllowedHeaders:   []string{"Origin", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           int(time.Duration(12 * time.Hour).Seconds()),
	})

	router.Use(middlewareFunc)

	v1Routes := router.Group("v1")
	{
		v1Routes.Use(middleware.KecyalokMiddleware())

		// Handle the POST requests at /v1/campaign
		v1Routes.POST(constants.Campaign, campaign.CreateCampaign)

		// Handle the PATCH requests at /v1/campaign
		v1Routes.PATCH(constants.Campaign, campaign.UpdateCampaign)

		// Handle the GET requests at /v1/campaign
		v1Routes.GET(constants.Campaign, campaign.GetCampaign)

		// Handle the GET requests at /v1/campaign/nearby
		v1Routes.GET(constants.CampaignNearby, nearby.GetCampaign)

		// Handle the PATCH requests at /v1/campaign/join
		v1Routes.PATCH(constants.CampaignJoin, join.JoinCampaign)

		// Handle the PATCH requests at /v1/campaign/leave
		v1Routes.PATCH(constants.CampaignLeave, join.LeaveCampaign)

	}

	unAuthRoutes := router.Group("v1")
	{
		// Handle the GET requests at /v1/statusNew
		unAuthRoutes.GET("/krakend.json", func(ctx *gin.Context) {
			lang := ctx.GetHeader("language")
			content, err := ioutil.ReadFile("utils/krakend/krakend.json")

			if err != nil {
				Msg := localization.GetMessage(lang, constants.InternalServerMessage, nil)
				utils.SendInternalServerError(ctx, Msg, "0", constants.IsJsonArray, nil)
				return
			}

			backendHost := applicationConfig.GetString(constants.ServerHost)

			krakendData := strings.ReplaceAll(string(content), "SERVER_HOST", backendHost)

			var result map[string]interface{}
			json.Unmarshal([]byte(krakendData), &result)

			// send success response
			ctx.JSON(http.StatusOK, result)
		})
	}

	return router
}
