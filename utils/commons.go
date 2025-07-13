package utils

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"

	"campaign-service/constants"

	genericModel "campaign-service/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

type ErrorRealm struct {
	GroupId             string `json:"groupId"`
	ClientId            string `json:"clientId"`
	UserId              string `json:"userId"`
	GuestUserId         string `json:"guestUserId"`
	ClientName          string `json:"clientName"`
	AccessToken         string `json:"accessToken"`
	InstitutePostgresId string `json:"institutePostgresId"`
	UserPostgresId      string `json:"userPostgresId"`
	DeleteRoles         bool   `json:"deleteRoles"`
	DeleteInstData      bool   `json:"deleteInstData"`
	NewUser             bool   `json:"newUser"`
}

type Enum interface {
	IsValid() bool
}

// check for empty string
func IsEmpty(feild string) bool {
	return feild == ""
}

func getData(respType int, data interface{}) interface{} {

	if data != nil {
		return data
	}
	type emptyData struct{}

	if respType == constants.IsJsonArray {
		return []emptyData{}

	} else if respType == constants.IsString {
		return ""
	}
	return emptyData{}
}

func SendBadRequest(ctx *gin.Context, msg, code string, respType int, err error) {
	data := getData(respType, nil)
	ctx.JSON(http.StatusBadRequest, genericModel.APIResponse{
		Status:    constants.APIRespErrorKey,
		Message:   msg,
		ErrorCode: code,
		Data:      data,
	})
}

func SendUnauthorized(ctx *gin.Context, msg, code string, respType int, err error) {
	ctx.JSON(http.StatusUnauthorized, genericModel.APIResponse{
		Status:    constants.APIRespErrorKey,
		Message:   msg,
		ErrorCode: code,
		Data:      getData(respType, nil),
	})
}

func SendRateLimit(ctx *gin.Context, msg, code string, respType int, err error) {
	ctx.JSON(http.StatusTooManyRequests, genericModel.APIResponse{
		Status:    constants.APIRespErrorKey,
		Message:   msg,
		ErrorCode: code,
		Data:      getData(respType, nil),
	})
}

func SendInternalServerError(ctx *gin.Context, msg, code string, respType int, err error) {
	data := getData(respType, nil)
	ctx.JSON(http.StatusInternalServerError, genericModel.APIResponse{
		Status:    constants.APIRespErrorKey,
		Message:   msg,
		ErrorCode: code,
		Data:      data,
	})
}

func SendNoDataFoundError(ctx *gin.Context, msg, code string, respType int, err error) {
	data := getData(respType, nil)
	ctx.JSON(http.StatusNotFound, genericModel.APIResponse{
		Status:    constants.APIRespErrorKey,
		Message:   msg,
		ErrorCode: code,
		Data:      data,
	})
}

func SendAccountDisable(ctx *gin.Context, msg, code string, respType int, err error) {
	ctx.JSON(http.StatusForbidden, genericModel.APIResponse{
		Status:    constants.APIRespErrorKey,
		Message:   msg,
		ErrorCode: code,
		Data:      getData(respType, nil),
	})
}

func SendConflict(ctx *gin.Context, code, msg string, respType int, data interface{}) {
	susdata := getData(respType, data)
	ctx.JSON(http.StatusConflict, genericModel.APIResponse{
		Status:    constants.ErrConflictMSG,
		Message:   msg,
		ErrorCode: code,
		Data:      susdata,
	})
}

func SendStatusOK(ctx *gin.Context, respType int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, genericModel.APIResponse{
		Status:  constants.APIRespSuccessKey,
		Message: msg,
		Data:    getData(respType, data),
	})
}

var (
	once sync.Once
	loc  *time.Location
)

func LogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	once.Do(func() {
		loc, _ = time.LoadLocation(constants.IST)
	})
	enc.AppendString(t.In(loc).Format(constants.TimeFormat))
}

func LevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}

func GetCurrentTimeInSeconds() int64 {
	return time.Now().Unix()
}

func GetCurrentTimeInMS() int64 {
	t := time.Now() //It will return time.Time object with current timestamp
	tUnixMilli := int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
	return tUnixMilli
}

func ValidateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(Enum)
	return value.IsValid()
}

func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func ErrorBasedOnResponse(ctx *gin.Context, msg string, respType int, err error) {

	switch err.Error() {

	//500
	default:
		SendInternalServerError(ctx, msg, constants.WLInternalServerErrorCode, respType, err)

	//400
	case errors.New(constants.BadRequestMessage).Error(), errors.New(constants.PasswordDoesNotMatchMessage).Error(),
		errors.New(constants.OldPasswordAndNewPasswordSameMessage).Error(), errors.New(constants.BusinessIdIsRequiredMessage).Error(), errors.New(constants.BrandIdIsRequiredMessage).Error(),
		errors.New(constants.InvalidStartDateMessage).Error(), errors.New(constants.InvalidEndDateMessage).Error(), errors.New(constants.EndDateBeforeStartDateMessage).Error(),
		errors.New(constants.MaxParticipantsLessThanMinParticipants).Error(), errors.New(constants.PriceMustBeGreaterThanZero).Error():
		SendBadRequest(ctx, msg, constants.WLBadRequestCode, respType, err)

	//401
	case errors.New(constants.UnauthorizedMessage).Error(), errors.New(constants.InvalidOldPasswordMessage).Error(), errors.New(constants.InvalidOTPMessage).Error():
		SendUnauthorized(ctx, msg, constants.WLUnauthorizedCode, respType, err)

	//403
	case errors.New(constants.UserNotActiveMessage).Error(), errors.New(constants.AccountDisableMessage).Error():
		SendAccountDisable(ctx, msg, constants.WLForbiddenCode, respType, err)

	//404
	case errors.New(constants.NotFoundMessage).Error(),
		errors.New(constants.AccountsDataNotFoundMessage).Error(), errors.New(constants.RoleNotFoundMessage).Error(), errors.New(constants.UserNotFoundMessage).Error(), errors.New(constants.CurrencyNotFoundMessage).Error(), errors.New(constants.SupportedLanguageNotFoundMessage).Error(), errors.New(constants.StateTaxNotFoundMessage).Error(), errors.New(constants.AccountTypeNotFoundMessage).Error(), errors.New(constants.CampaignNotFoundMessage).Error():
		SendNoDataFoundError(ctx, msg, constants.WLNoDataFoundCode, respType, err)

	//409
	case errors.New(constants.ConflictMessage).Error(), errors.New(constants.UserAlreadyExistsMessage).Error(), errors.New(constants.CampaignNameAlreadyExistsMessage).Error():
		SendConflict(ctx, constants.WLDataConflictCode, msg, respType, nil)

	//429
	case errors.New(constants.TooManyRequestsMessage).Error():
		SendRateLimit(ctx, msg, constants.WLRateLimitCode, respType, err)

	}

}

func GetSkipLimit(ctx *gin.Context) (int, int) {

	limitQuery := ctx.Request.URL.Query().Get("limit")
	skipQuery := ctx.Request.URL.Query().Get("skip")

	limit, _ := StringToInt(limitQuery)
	skip, _ := StringToInt(skipQuery)

	if limit == 0 || limit > 20 {
		limit = 20
	}

	return skip, limit
}

func RemoveSpaces(input string) string {
	// Replace all spaces with empty strings
	result := strings.ReplaceAll(input, " ", "")
	return result
}
