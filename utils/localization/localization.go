package localization

import (
	"campaign-service/constants"
	"campaign-service/logger"
	"campaign-service/utils/configs"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle
var bundleOnce sync.Once
var supportedLanguages []string

func GetLocalizer(ctx context.Context) *i18n.Localizer {
	ensureBundle("")
	return i18n.NewLocalizer(Bundle, constants.DefaultLanguage)
}

func GetMessageWithoutTemplate(localizer *i18n.Localizer, id string) string {
	message, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
	})
	return message
}

func ensureBundle(path string) {
	bundleOnce.Do(func() {
		log := logger.GetLoggerWithoutContext()
		defaultLang := constants.DefaultLanguage

		supportedLanguages = []string{defaultLang}

		LanguageConfig, err := configs.Get(constants.LanguageConfig)
		if err != nil {
			log.With(zap.Error(err)).Warn("failed to load language config; falling back to default language")
		} else if LanguageConfig != nil {
			langs := LanguageConfig.GetStringSlice(constants.LanguageListKey)
			if len(langs) > 0 {
				// Deduplicate + trim and ensure default exists.
				seen := make(map[string]struct{}, len(langs)+1)
				res := make([]string, 0, len(langs)+1)
				add := func(l string) {
					l = strings.TrimSpace(l)
					if l == "" {
						return
					}
					if _, ok := seen[l]; ok {
						return
					}
					seen[l] = struct{}{}
					res = append(res, l)
				}

				add(defaultLang)
				for _, l := range langs {
					add(l)
				}
				supportedLanguages = res
			}
		}

		Bundle = i18n.NewBundle(language.Make(defaultLang))
		Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		for _, lang := range supportedLanguages {
			var filePattern string
			if path != "" {
				filePattern = path + "locales/%v.json"
			} else {
				filePattern = "locales/%v.json"
			}

			file := fmt.Sprintf(filePattern, lang)
			if _, err := Bundle.LoadMessageFile(file); err != nil {
				// Missing locale files should not crash the service.
				log.With(zap.Error(err), zap.String("lang", lang), zap.String("file", file)).
					Warn("failed to load locale file")
			}
		}
	})
}

func LoadBundle(path string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ensureBundle(path)

		lan := constants.DefaultLanguage
		if headerLang := strings.TrimSpace(ctx.GetHeader(constants.LanguageString)); headerLang != "" {
			if pos(supportedLanguages, headerLang) != -1 {
				lan = headerLang
			}
		}
		ctx.Set(constants.LanguageString, lan)
	}
}

func GetMessage(lang interface{}, id string, templateData interface{}) string {
	ensureBundle("")

	selectedLang := constants.DefaultLanguage
	if lang != nil {
		if s, ok := lang.(string); ok {
			s = strings.TrimSpace(s)
			if s != "" {
				selectedLang = s
			}
		}
	}

	if pos(supportedLanguages, selectedLang) == -1 {
		selectedLang = constants.DefaultLanguage
	}

	if Bundle == nil {
		return id
	}

	localizer := i18n.NewLocalizer(Bundle, selectedLang)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
		TemplateData: templateData,
	})
	if err != nil || message == "" {
		return id
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
