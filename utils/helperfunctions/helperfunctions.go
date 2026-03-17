package helperfunctions

import (
	"campaign-service/constants"
	"campaign-service/library/mongoDb"
	"campaign-service/library/postgres"
	"campaign-service/library/redis_provider"
	"campaign-service/logger"
	"campaign-service/models"
	"campaign-service/utils"
	"campaign-service/utils/localization"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// validate the user if it exists
	func ValidateUserExists(ctx context.Context, userID string) (models.Users, error) {
		log := logger.GetLoggerWithoutContext()

		// get the user collection
		userCollection := mongoDb.GetCollection(constants.MongoUserCollection)

		// convert string to object id
		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Error("Failed to convert userID to objectID", zap.Error(err))
			return models.Users{}, err
		}

		// find if the user exists in the database
		var user models.Users
		err = userCollection.FindOne(ctx, bson.M{"_id": userObjID, "status": "Active"}).Decode(&user)
		if err != nil {
			log.Error("Failed to get user by ID", zap.Error(err))
			return models.Users{}, err
		}

		return user, nil
	}

	// InvalidateAllCampaignUserCache invalidates all campaign-related cache for a specific user
	func InvalidateAllCampaignUserCache(ctx context.Context, userID string) error {
		log := logger.GetLoggerWithoutContext()
		redis := redis_provider.Client

		// Delete individual campaign keys
		indexKey := fmt.Sprintf("campaign:user:%s:index", userID)
		campaignIDs, err := redis.SMembers(ctx, indexKey).Result()
		if err != nil && err.Error() != "redis: nil" {
			log.Error("Failed to get campaign IDs from redis set", zap.Error(err))
			return fmt.Errorf("failed to get campaign IDs from redis set: %w", err)
		}

		// Delete individual campaign keys
		var campaignKeys []string
		for _, id := range campaignIDs {
			campaignKeys = append(campaignKeys, fmt.Sprintf("campaign:user:%s:%s", userID, id))
		}

		if len(campaignKeys) > 0 {
			if err := redis.Del(ctx, campaignKeys...).Err(); err != nil {
				log.Error("Failed to delete campaign keys", zap.Error(err))
				return fmt.Errorf("failed to delete campaign keys: %w", err)
			}
		}

		// Delete the index key
		if err := redis.Del(ctx, indexKey).Err(); err != nil {
			log.Error("Failed to delete index key", zap.Error(err))
			return fmt.Errorf("failed to delete index key: %w", err)
		}

		// Delete all query cache keys for this user
		queryPattern := fmt.Sprintf("campaign:user:%s:*", userID)
		queryCacheKeys, err := redis.Keys(ctx, queryPattern).Result()
		if err != nil {
			log.Error("Failed to get query cache keys", zap.Error(err))
			return fmt.Errorf("failed to get query cache keys: %w", err)
		}

		// Filter out individual campaign keys and index key (already deleted)
		var filteredQueryKeys []string
		for _, key := range queryCacheKeys {
			// Skip individual campaign keys (they have UUID format at the end)
			// Skip index key (already deleted)
			if !strings.Contains(key, ":index") && !isIndividualCampaignKey(key) {
				filteredQueryKeys = append(filteredQueryKeys, key)
			}
		}

		if len(filteredQueryKeys) > 0 {
			if err := redis.Del(ctx, filteredQueryKeys...).Err(); err != nil {
				log.Error("Failed to delete query cache keys", zap.Error(err))
				return fmt.Errorf("failed to delete query cache keys: %w", err)
			}
		}

		log.Info("Successfully invalidated all campaign cache for user", zap.String("userID", userID))
		return nil
	}

	// isIndividualCampaignKey checks if a key is an individual campaign key
	func isIndividualCampaignKey(key string) bool {
		parts := strings.Split(key, ":")
		if len(parts) < 4 {
			return false
		}
		// Individual campaign keys have format: campaign:user:userID:campaignID
		// campaignID is typically a UUID, so we check if the last part looks like a UUID
		lastPart := parts[len(parts)-1]
		return len(lastPart) == 36 && strings.Count(lastPart, "-") == 4
	}

	// GetActiveCampaignsFromRedis gets all active campaigns from Redis
	// Note: This function should be used carefully as it can be expensive for large datasets
	func GetActiveCampaignsFromRedis(ctx context.Context) ([]models.Campaign, error) {
		log := logger.GetLoggerWithoutContext()
		redis := redis_provider.Client

		// Use a more specific pattern to avoid getting query cache keys
		cachePattern := "campaign:user:*:*"
		keys, err := redis.Keys(ctx, cachePattern).Result()
		if err != nil {
			log.Error("Failed to get campaign keys from redis", zap.Error(err))
			return nil, fmt.Errorf("failed to get campaign keys from redis: %w", err)
		}

		var campaigns []models.Campaign
		var validKeys []string

		// Filter to only individual campaign keys
		for _, key := range keys {
			if isIndividualCampaignKey(key) && !strings.Contains(key, ":index") {
				validKeys = append(validKeys, key)
			}
		}

		// Process keys in batches to avoid memory issues
		batchSize := 100
		for i := 0; i < len(validKeys); i += batchSize {
			end := i + batchSize
			if end > len(validKeys) {
				end = len(validKeys)
			}

			batch := validKeys[i:end]
			results, err := redis.MGet(ctx, batch...).Result()
			if err != nil {
				log.Error("Failed to get campaign batch from redis", zap.Error(err))
				continue
			}

			for _, result := range results {
				if result == nil {
					continue
				}
				var campaign models.Campaign
				if err := json.Unmarshal([]byte(result.(string)), &campaign); err == nil {
					campaigns = append(campaigns, campaign)
				}
			}
		}

		log.Info("Retrieved campaigns from redis", zap.Int("count", len(campaigns)))
		return campaigns, nil
	}

	// AddCampaignToUserIndex adds a campaign ID to the user's campaign index set
	func AddCampaignToUserIndex(ctx context.Context, userID, campaignID string) error {
		log := logger.GetLoggerWithoutContext()
		redis := redis_provider.Client

		indexKey := fmt.Sprintf("campaign:user:%s:index", userID)

		// Add to set
		if err := redis.SAdd(ctx, indexKey, campaignID).Err(); err != nil {
			log.Error("Failed to add campaign to user index",
				zap.String("userID", userID),
				zap.String("campaignID", campaignID),
				zap.Error(err))
			return fmt.Errorf("failed to add campaign to user index: %w", err)
		}

		// Set TTL for the index set
		if err := redis.Expire(ctx, indexKey, 30*time.Minute).Err(); err != nil {
			log.Error("Failed to set TTL for user index set",
				zap.String("userID", userID),
				zap.Error(err))
			return fmt.Errorf("failed to set TTL for user index set: %w", err)
		}

		return nil
	}

	// GetCampaignFromRedis retrieves a specific campaign from Redis
	func GetCampaignFromRedis(ctx context.Context, userID, campaignID string) (*models.Campaign, error) {
		log := logger.GetLoggerWithoutContext()
		redis := redis_provider.Client

		cacheKey := fmt.Sprintf("campaign:user:%s:%s", userID, campaignID)
		val, err := redis.Get(ctx, cacheKey).Result()
		if err != nil {
			if err.Error() == "redis: nil" {
				return nil, fmt.Errorf("campaign not found in cache")
			}
			log.Error("Failed to get campaign from redis",
				zap.String("userID", userID),
				zap.String("campaignID", campaignID),
				zap.Error(err))
			return nil, fmt.Errorf("failed to get campaign from redis: %w", err)
		}

		var campaign models.Campaign
		if err := json.Unmarshal([]byte(val), &campaign); err != nil {
			log.Error("Failed to unmarshal campaign from redis",
				zap.String("userID", userID),
				zap.String("campaignID", campaignID),
				zap.Error(err))
			return nil, fmt.Errorf("failed to unmarshal campaign: %w", err)
		}

		return &campaign, nil
	}

	// SetCampaignInRedis caches a campaign in Redis
	func SetCampaignInRedis(ctx context.Context, userID string, campaign *models.Campaign) error {
		log := logger.GetLoggerWithoutContext()
		redis := redis_provider.Client

		cacheKey := fmt.Sprintf("campaign:user:%s:%s", userID, campaign.ID.String())

		campaignJSON, err := json.Marshal(campaign)
		if err != nil {
			log.Error("Failed to marshal campaign for caching",
				zap.String("userID", userID),
				zap.String("campaignID", campaign.ID.String()),
				zap.Error(err))
			return fmt.Errorf("failed to marshal campaign: %w", err)
		}

		// Set individual campaign with 30-minute TTL
		if err := redis.Set(ctx, cacheKey, campaignJSON, 30*time.Minute).Err(); err != nil {
			log.Error("Failed to set campaign in redis",
				zap.String("userID", userID),
				zap.String("campaignID", campaign.ID.String()),
				zap.Error(err))
			return fmt.Errorf("failed to set campaign in redis: %w", err)
		}

		// Add to user's campaign index
		if err := AddCampaignToUserIndex(ctx, userID, campaign.ID.String()); err != nil {
			log.Error("Failed to add campaign to user index",
				zap.String("userID", userID),
				zap.String("campaignID", campaign.ID.String()),
				zap.Error(err))
			// Don't return error here as the main caching succeeded
		}

		return nil
	}

	// RemoveCampaignFromRedis removes a specific campaign from Redis
	func RemoveCampaignFromRedis(ctx context.Context, userID, campaignID string) error {
		log := logger.GetLoggerWithoutContext()
		redis := redis_provider.Client

		cacheKey := fmt.Sprintf("campaign:user:%s:%s", userID, campaignID)
		indexKey := fmt.Sprintf("campaign:user:%s:index", userID)

		// Remove individual campaign
		if err := redis.Del(ctx, cacheKey).Err(); err != nil {
			log.Error("Failed to delete campaign from redis",
				zap.String("userID", userID),
				zap.String("campaignID", campaignID),
				zap.Error(err))
			return fmt.Errorf("failed to delete campaign from redis: %w", err)
		}

		// Remove from user's campaign index
		if err := redis.SRem(ctx, indexKey, campaignID).Err(); err != nil {
			log.Error("Failed to remove campaign from user index",
				zap.String("userID", userID),
				zap.String("campaignID", campaignID),
				zap.Error(err))
			return fmt.Errorf("failed to remove campaign from user index: %w", err)
		}

		return nil
	}

	// InvalidateQueryCache invalidates all query cache for a specific user
	func InvalidateQueryCache(ctx context.Context, userID string) error {
		log := logger.GetLoggerWithoutContext()
		redis := redis_provider.Client

		// Delete all query cache keys for this user
		queryPattern := fmt.Sprintf("campaign:user:%s:*", userID)
		queryCacheKeys, err := redis.Keys(ctx, queryPattern).Result()
		if err != nil {
			log.Error("Failed to get query cache keys", zap.Error(err))
			return fmt.Errorf("failed to get query cache keys: %w", err)
		}

		// Filter out individual campaign keys and index key
		var filteredQueryKeys []string
		for _, key := range queryCacheKeys {
			if !strings.Contains(key, ":index") && !isIndividualCampaignKey(key) {
				filteredQueryKeys = append(filteredQueryKeys, key)
			}
		}

		if len(filteredQueryKeys) > 0 {
			if err := redis.Del(ctx, filteredQueryKeys...).Err(); err != nil {
				log.Error("Failed to delete query cache keys", zap.Error(err))
				return fmt.Errorf("failed to delete query cache keys: %w", err)
			}
		}

		log.Info("Successfully invalidated query cache for user", zap.String("userID", userID))
		return nil
	}

	// AddCampaignToSpatialIndex adds a campaign to the spatial index
	func AddCampaignToSpatialIndex(ctx context.Context, campaignID string, latitude, longitude float64) error {
		return redis_provider.AddCampaignToSpatialIndex(ctx, campaignID, latitude, longitude)
	}

	// RemoveCampaignFromSpatialIndex removes a campaign from the spatial index
	func RemoveCampaignFromSpatialIndex(ctx context.Context, campaignID string) error {
		return redis_provider.RemoveCampaignFromSpatialIndex(ctx, campaignID)
	}

	// GetCampaignsInRadius gets campaign IDs within a specified radius
	func GetCampaignsInRadius(ctx context.Context, latitude, longitude, radius float64, unit string) ([]string, error) {
		return redis_provider.GetCampaignsInRadius(ctx, latitude, longitude, radius, unit)
	}
