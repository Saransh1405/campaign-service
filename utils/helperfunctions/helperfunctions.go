package helperfunctions

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"users-service/constants"
	"users-service/library/postgres"
	"users-service/logger"
	"users-service/models"
	"users-service/utils"
	"users-service/utils/localization"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GeneratePassword() string {
	rand.Seed(time.Now().Unix())
	lowerCharSet := constants.ABCDLower
	upperCharSet := constants.ABCDUpper
	specialCharSet := constants.SpecialCharSet2
	numberSet := constants.Number
	allCharSet := lowerCharSet + upperCharSet + specialCharSet + numberSet
	minSpecialChar := 2
	minNum := 2
	minUpperCase := 2
	passwordLength := 13

	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func InsertLogsIntoStatusLogs(status models.Status, notes, actionByUserRole, actionByUserId string, entityId uuid.UUID) {

	var statusLogs models.StatusLogs

	statusLogs.Status = status
	statusLogs.ActionByUserRole = actionByUserRole
	statusLogs.ActionByUserId = actionByUserId
	statusLogs.EntityId = entityId
	statusLogs.Notes = notes
	statusLogs.Timestamp = time.Now()

	postgres.DB.Create(&statusLogs)

}

func ValidateRequestData(ctx *gin.Context, request interface{}, b binding.Binding) error {

	lang, _ := ctx.Get(constants.LanguageString)
	log := logger.GetLogger(ctx)

	err := ctx.ShouldBindWith(request, b)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
		var verr validator.ValidationErrors
		fields := []string{}
		if errors.As(err, &verr) {
			for _, f := range verr {
				fields = append(fields, f.Field())
			}
		}
		Badrequestmsg := localization.GetMessage(lang, constants.BadRequestMessage, map[string]interface{}{
			"Fields": strings.Join(fields, ", "),
		})
		utils.SendBadRequest(ctx, constants.BadRequestErr, Badrequestmsg, constants.IsString, err)
		return err
	}

	return nil
}

func ValidateRequestDataParam(ctx *gin.Context, request interface{}) error {

	lang, _ := ctx.Get(constants.LanguageString)
	log := logger.GetLogger(ctx)

	err := ctx.ShouldBind(request)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
		var verr validator.ValidationErrors
		fields := []string{}
		if errors.As(err, &verr) {
			for _, f := range verr {
				fields = append(fields, f.Field())
			}
		}
		Badrequestmsg := localization.GetMessage(lang, constants.BadRequestMessage, map[string]interface{}{
			"Fields": strings.Join(fields, ", "),
		})
		utils.SendBadRequest(ctx, constants.BadRequestErr, Badrequestmsg, constants.IsString, err)
		return err
	}

	return nil
}

func GenerateID() string {
	rand.Seed(time.Now().UTC().UnixNano())
	lowerCharSet := constants.ABCDLower
	numberSet := constants.Number
	allCharSet := lowerCharSet + "-" + numberSet
	minNum := 7
	IDLength := 20

	var password strings.Builder

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	remainingLength := IDLength - minNum
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func CheckUserAlreadyExistsForProperty(countryCode, phone, email string) (bool, error) {

	var foundUser models.Users

	find := postgres.DB.Where("country_code = ? AND phone = ? AND email = ?", countryCode, phone, email).First(&foundUser)

	if find.RowsAffected == 0 {
		return false, nil
	}

	if find.Error != nil {
		return false, find.Error
	}

	return true, nil
}

func AddLogs(trigger, enitity, enitityId, clientName, actionById string, oldData, newData interface{}) {

	insertLogs := models.Logs{
		Trigger:    trigger,
		Entity:     enitity,
		EntityId:   enitityId,
		ClientName: clientName,
		ActionById: actionById,
		OldData:    oldData,
		NewData:    newData,
		Timestamp:  time.Now(),
	}

	insert := postgres.DB.Create(&insertLogs)

	if insert.Error != nil {
		fmt.Printf("insert.Error: %v\n", insert.Error)
	}

}
