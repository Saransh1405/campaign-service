package api

import (
	"campaign-service/api/campaign"
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

		// // Handle the POST requests at /v1/role
		// v1Routes.POST(constants.Role, role.Post)

		// // Handle the PATCH requests at /v1/role/:id
		// v1Routes.PATCH(constants.RoleWithId, role.Patch)

		// // Handle the GET requests at /v1/role
		// v1Routes.GET(constants.Role, role.GetAll)

		// // Handle the GET requests at /v1/role/:id
		// v1Routes.GET(constants.RoleWithId, role.Get)

		// // Handle the DELETE requests at /v1/role
		// v1Routes.DELETE(constants.Role, role.Delete)

		// // Handle the POST requests at /v1/verifyEmail
		// v1Routes.POST(constants.VerifyEmail, verifyEmail.Post)

		// // Handle the POST requests at /v1/sendVerifyEmail
		// v1Routes.POST(constants.SendVerifyEmail, verifyEmail.PostSendEmail)

		// // Handle the PATCH requests at /v1/verifyEmail
		// v1Routes.PATCH(constants.VerifyEmail, verifyEmail.Patch)

		// // Handle the POST requests at /v1/verifyPhone
		// v1Routes.POST(constants.VerifyPhone, verifyPhone.Post)

		// // Handle the PATCH requests at /v1/changeEmail
		// v1Routes.PATCH(constants.ChangeEmail, changeEmail.Patch)

		// // Handle the PATCH requests at /v1/changePhone
		// v1Routes.PATCH(constants.ChangePhone, changePhone.Patch)

		// // Handle the PATCH requests at /v1/password
		// v1Routes.PATCH(constants.Password, password.Patch)

		// // Handle the POST requests at /v1/resetPassword
		// v1Routes.POST(constants.ResetPassword, password.Post)

		// // Handle the POST requests at /v1/user
		// v1Routes.POST(constants.User, users.Post)

		// // Handle the PATCH requests at /v1/user/:id
		// v1Routes.PATCH(constants.UserWithId, users.Patch)

		// // Handle the GET requests at /v1/user/:id
		// v1Routes.GET(constants.UserWithId, users.Get)

		// // Handle the GET requests at /v1/users
		// v1Routes.GET(constants.User, users.GetAll)

		// // Handle the DELETE requests at /v1/user/:id
		// v1Routes.DELETE(constants.User, users.Delete)

		// //Handel the POST requests at /v1/countries
		// v1Routes.POST(constants.Countries, countries.Post)

		// //Handel the GET requests at /v1/countries
		// v1Routes.GET(constants.Countries, countries.GetAll)

		// //Handel the GET requests at /v1/countries/:id
		// v1Routes.GET(constants.CountriesWithId, countries.Get)

		// //Handel the GET requests at /v1/countryLevelTax/:countryId
		// v1Routes.GET(constants.CountryLevelTaxWithCountryId, countries.GetCountryLevelTax)

		// //Handel the PATCH requests at /v1/countryLevelTax/:id
		// v1Routes.PATCH(constants.CountryLevelTaxWithId, countries.PatchCountryLevelTax)

		// //Handel the DELETE requests at /v1/countryLevelTax
		// v1Routes.DELETE(constants.CountryLevelTax, countries.DeleteCountryLevelTax)

		// // Handle the GET requests at /v1/propertyTypes
		// v1Routes.GET(constants.PropertyTypes, propertyTypes.GetAll)

		// // Handle the GET requests at /v1/propertyTypes/:id
		// v1Routes.GET(constants.PropertyTypesWithId, propertyTypes.Get)

		// // Handle the POST requests at /v1/propertyTypes
		// v1Routes.POST(constants.PropertyTypes, propertyTypes.Post)

		// // Handle the PATCH requests at /v1/propertyTypes/:id
		// v1Routes.PATCH(constants.PropertyTypesWithId, propertyTypes.Patch)

		// // Handle the DELETE requests at /v1/propertyTypes
		// v1Routes.DELETE(constants.PropertyTypes, propertyTypes.Delete)

		// // Handle the GET requests at /v1/propertyAmenities
		// v1Routes.GET(constants.PropertyAmenities, propertyAmenities.GetAll)

		// // Handle the GET requests at /v1/propertyAmenities/:id
		// v1Routes.GET(constants.PropertyAmenitiesWithId, propertyAmenities.Get)

		// // Handle the POST requests at /v1/propertyAmenities
		// v1Routes.POST(constants.PropertyAmenities, propertyAmenities.Post)

		// // Handle the PATCH requests at /v1/propertyAmenities/:id
		// v1Routes.PATCH(constants.PropertyAmenitiesWithId, propertyAmenities.Patch)

		// // Handle the DELETE requests at /v1/propertyAmenities
		// v1Routes.DELETE(constants.PropertyAmenities, propertyAmenities.Delete)

		// // Handle the POST requests at /v1/currencies
		// v1Routes.POST(constants.Currencies, currencies.Post)

		// // Handle the PATCH requests at /v1/currencies/:id
		// v1Routes.PATCH(constants.CurrenciesWithId, currencies.Patch)

		// // Handle the GET requests at /v1/currencies/:id
		// v1Routes.GET(constants.CurrenciesWithId, currencies.Get)

		// // Handle the GET requests at /v1/currencies
		// v1Routes.GET(constants.Currencies, currencies.GetAll)

		// // Handle the DELETE requests at /v1/currencies
		// v1Routes.DELETE(constants.Currencies, currencies.Delete)

		// //Handle the POST requests at /v1/roomViews
		// v1Routes.POST(constants.RoomViews, roomViews.Post)

		// //Handle the PATCH requests at /v1/roomViews/:id
		// v1Routes.PATCH(constants.RoomViewsWithId, roomViews.Patch)

		// // Handle the GET requests at /v1/roomViews/:id
		// v1Routes.GET(constants.RoomViewsWithId, roomViews.Get)

		// //Handle the GET requests at /v1/roomViews
		// v1Routes.GET(constants.RoomViews, roomViews.GetAll)

		// // Handle the DELETE requests at /v1/roomViews
		// v1Routes.DELETE(constants.RoomViews, roomViews.Delete)

		// // Handle the GET requests at /v1/documents
		// v1Routes.GET(constants.Documents, documents.GetAll)

		// // Handle the GET requests at /v1/document/:id
		// v1Routes.GET(constants.DocumentWithId, documents.Get)

		// // Handle the Post requests at /v1/document
		// v1Routes.POST(constants.Document, documents.Post)

		// // Handle the PATCH requests at /v1/document/:id
		// v1Routes.PATCH(constants.DocumentWithId, documents.Patch)

		// // Handle the GET requests at /v1/propertyModules/:id
		// v1Routes.GET(constants.PropertyModulesWithId, propertyModules.Get)

		// // Handle the PATCH requests at /v1/propertyModules/:id
		// v1Routes.PATCH(constants.PropertyModulesWithId, propertyModules.Patch)

		// // Handle the GET requests at /v1/businessMetaData/:id
		// v1Routes.GET(constants.BusinessMetaDataWithId, businessMetaData.Get)

		// // Handle the GET requests at /v1/businessMetaData
		// v1Routes.GET(constants.BusinessMetaData, businessMetaData.GetAll)

		// // Handle the PATCH requests at /v1/businessMetaData/:id
		// v1Routes.PATCH(constants.BusinessMetaDataWithId, businessMetaData.Patch)

		// // Handle the POST requests at /v1/businessMetaData
		// v1Routes.POST(constants.BusinessMetaData, businessMetaData.Post)

		// // Handle the GET requests at /v1/reasons
		// v1Routes.GET(constants.Reasons, reasons.GetAll)

		// // Handle the GET requests at /v1/reasons/:id
		// v1Routes.GET(constants.ReasonsWithId, reasons.Get)

		// // Handle the POST requests at /v1/reasons
		// v1Routes.POST(constants.Reasons, reasons.Post)

		// // Handle the PATCH requests at /v1/reasons/:id
		// v1Routes.PATCH(constants.ReasonsWithId, reasons.Patch)

		// // Handle the DELETE requests at /v1/reasons
		// v1Routes.DELETE(constants.Reasons, reasons.Delete)

		// // Handle the POST requests at /v1/supportedLanguages
		// v1Routes.POST(constants.SupportedLanguages, supportedLanguages.Post)

		// // Handle the GET requests at /v1/supportedLanguages/:id
		// v1Routes.GET(constants.SupportedLanguagesWithId, supportedLanguages.Get)

		// // Handle the GET requests at /v1/supportedLanguages
		// v1Routes.GET(constants.SupportedLanguages, supportedLanguages.GetAll)

		// // Handle the PATCH requests at /v1/supportedLanguages/:id
		// v1Routes.PATCH(constants.SupportedLanguagesWithId, supportedLanguages.Patch)

		// // Handle the DELETE requests at /v1/supportedLanguages
		// v1Routes.DELETE(constants.SupportedLanguages, supportedLanguages.Delete)

		// // Handle the GET requests at /v1/products/:id
		// v1Routes.GET(constants.ProductsWithId, products.Get)

		// // Handle the GET requests at /v1/products
		// v1Routes.GET(constants.Products, products.GetAll)

		// // Handle the POST requests at /v1/products
		// v1Routes.POST(constants.Products, products.Post)

		// // Handle the PATCH requests at /v1/products/:id
		// v1Routes.PATCH(constants.ProductsWithId, products.Patch)

		// // Handle the DELETE requests at /v1/products/:id
		// v1Routes.DELETE(constants.Products, products.Delete)

		// // Handle the GET requests at /v1/taxFields/:countryId
		// v1Routes.GET(constants.TaxFieldsWithCountryId, taxFields.GetAll)

		// // Handle the GET requests at /v1/bankFields/:countryId
		// v1Routes.GET(constants.BankFieldsWithCountryId, bankFields.GetAll)

		// // Handle the GET requests at /v1/properties/:id
		// v1Routes.GET(constants.PropertiesWithId, properties.Get)

		// // Handle the GET requests at /v1/properties
		// v1Routes.GET(constants.Properties, properties.GetAll)

		// // Handle the POST requests at /v1/properties
		// v1Routes.POST(constants.Properties, properties.Post)

		// // Handle the PATCH requests at /v1/properties/:id
		// v1Routes.PATCH(constants.PropertiesWithId, properties.Patch)

		// // Handle the DELETE requests at /v1/properties
		// v1Routes.DELETE(constants.Properties, properties.Delete)

		// // Handle the GET requests at /v1/brands/:id
		// v1Routes.GET(constants.BrandsWithId, brands.Get)

		// // Handle the GET requests at /v1/brands
		// v1Routes.GET(constants.Brands, brands.GetAll)

		// // Handle the POST requests at /v1/brands
		// v1Routes.POST(constants.Brands, brands.Post)

		// // Handle the PATCH requests at /v1/brands/:id
		// v1Routes.PATCH(constants.BrandsWithId, brands.Patch)

		// // Handle the DELETE requests at /v1/brands
		// v1Routes.DELETE(constants.Brands, brands.Delete)

		// // Handle the GET requests at /v1/businesses/:id
		// v1Routes.GET(constants.BusinessesWithId, businesses.Get)

		// // Handle the GET requests at /v1/businesses
		// v1Routes.GET(constants.Businesses, businesses.GetAll)

		// // Handle the POST requests at /v1/businesses
		// v1Routes.POST(constants.Businesses, businesses.Post)

		// // Handle the PATCH requests at /v1/businesses/:id
		// v1Routes.PATCH(constants.BusinessesWithId, businesses.Patch)

		// // Handle the DELETE requests at /v1/businesses
		// v1Routes.DELETE(constants.Businesses, businesses.Delete)

		// // Handle the GET requests at /v1/stateTaxes/:id
		// v1Routes.GET(constants.StateTaxesWithId, stateTaxes.Get)

		// // Handle the GET requests at /v1/stateTaxes
		// v1Routes.GET(constants.StateTaxes, stateTaxes.GetAll)

		// // Handle the POST requests at /v1/stateTaxes
		// v1Routes.POST(constants.StateTaxes, stateTaxes.Post)

		// // Handle the PATCH requests at /v1/stateTaxes/:id
		// v1Routes.PATCH(constants.StateTaxesWithId, stateTaxes.Patch)

		// // Handle the DELETE requests at /v1/stateTaxes
		// v1Routes.DELETE(constants.StateTaxes, stateTaxes.Delete)

		// // Handle the POST requests at /v1/holidays
		// v1Routes.POST(constants.Holidays, holidays.Post)

		// // Handle the PATCH requests at /v1/holidays/:id
		// v1Routes.PATCH(constants.HolidaysWithId, holidays.Patch)

		// // Handle the GET requests at /v1/holidays/:id
		// v1Routes.GET(constants.HolidaysWithId, holidays.Get)

		// // Handle the GET requests at /v1/holidays
		// v1Routes.GET(constants.Holidays, holidays.GetAll)

		// // Handle the DELETE requests at /v1/holidays
		// v1Routes.DELETE(constants.Holidays, holidays.Delete)

		// // Handle the PATCH request at /v1/workHours/:accountId
		// v1Routes.PATCH(constants.WorkHoursWithAccountId, workingHours.Patch)

		// // Handle the POST requests at /v1/leaveTypes
		// v1Routes.POST(constants.LeaveTypes, leaveTypes.Post)

		// // Handle the PATCH requests at /v1/leaveTypes/:id
		// v1Routes.PATCH(constants.LeaveTypesWithId, leaveTypes.Patch)

		// // Handle the GET requests at /v1/leaveTypes/:id
		// v1Routes.GET(constants.LeaveTypesWithId, leaveTypes.Get)

		// // Handle the GET requests at /v1/leaveTypes
		// v1Routes.GET(constants.LeaveTypes, leaveTypes.GetAll)

		// // Handle the DELETE requests at /v1/leaveTypes
		// v1Routes.DELETE(constants.LeaveTypes, leaveTypes.Delete)

		// // Handle the PATCH requests at /v1/documentsForEnlisting/:accountId
		// v1Routes.PATCH(constants.DocumentsForEnlistingWithAccountId, documentsForEnlisting.Patch)

	}

	unAuthRoutes := router.Group("v1")
	{
		// // Handle the POST requests at /v1/login
		// unAuthRoutes.POST(constants.Login, login.Post)

		// // Handle the POST requests at /v1/logout
		// unAuthRoutes.POST(constants.Logout, logout.Post)

		// // Handle the POST requests at /v1/loginWithOtp
		// unAuthRoutes.POST(constants.LoginWithOtp, otpLogin.Post)

		// // Handle the PATCH requests at /v1/loginWithOtp
		// unAuthRoutes.PATCH(constants.LoginWithOtp, otpLogin.Patch)

		// // Handle the POST requests at /v1/signup
		// unAuthRoutes.POST(constants.Signup, signup.Post)

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
