package localization

import (
	"campaign-service/constants"
	"campaign-service/logger"
	"campaign-service/utils/configs"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// Bundle - to log output on console.
var Bundle *i18n.Bundle

func GetLocalizer(ctx context.Context) *i18n.Localizer {
	log := logger.GetLoggerWithoutContext()
	lan := constants.DefaultLanguage
	Bundle = i18n.NewBundle(language.Make(lan))
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err := Bundle.LoadMessageFile(fmt.Sprintf(constants.LanguageJsonFilePath, lan))
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}
	return i18n.NewLocalizer(Bundle, lan)
}

func GetMessageWithoutTemplate(localizer *i18n.Localizer, id string) string {
	message, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
	})
	return message
}

// LoadBundle local locales file.
func LoadBundle(path string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Initialize i18n
		lan := constants.DefaultLanguage
		Bundle = i18n.NewBundle(language.Make(lan))
		Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		LanguageConfig, _ := configs.Get(constants.LanguageConfig)
		LanguageList := LanguageConfig.GetStringSlice(constants.LanguageListKey)

		LL := make([]string, len(LanguageList))
		for i, v := range LanguageList {
			LL[i] = fmt.Sprint(v)
		}
		log := logger.GetLoggerWithoutContext()
		log.Info(fmt.Sprintf("Language List %v", LL))
		for _, lang := range LL {
			if path != "" {
				Bundle.MustLoadMessageFile(fmt.Sprintf(path+"locales/%v.json", lang))
			} else {
				Bundle.MustLoadMessageFile(fmt.Sprintf("locales/%v.json", lang))
			}

		}
		ctx.Set("language", lan)
	}
}

// GetMessage get message from local file.
func GetMessage(lang interface{}, id string, templateData interface{}) string {

	log := logger.GetLoggerWithoutContext()
	language := constants.DefaultLanguage
	if lang != nil {
		language = lang.(string)
	}
	languageString := []string{}
	LanguageConfig, err := configs.Get(constants.LanguageConfig)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}

	languageString = LanguageConfig.GetStringSlice(constants.LanguageListKey)

	if pos(languageString, language) == -1 {
		language = constants.DefaultLanguage
	}

	localizer := i18n.NewLocalizer(Bundle, language)

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
		TemplateData: templateData,
	})
	if err != nil || message == "" {
		message = id
	}
	return message
}

func pos(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
