package models

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;index" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `gorm:"index" json:"deletedAt"`
} //@name Model

type APIResponse struct {
	Status    string      `json:"status" example:"success"`
	Message   string      `json:"message,omitempty" example:"success"`
	ErrorCode string      `json:"errorCode,omitempty" example:"0"`
	Data      interface{} `json:"data,omitempty"`
} //@name APIResponse

// ErrorResponse is the common error response body to be used in case of any error
type ErrorResponse struct {
	Description string `json:"description" example:"error description"`
} //@name ErrorResponse

// Data from API Gateway
type UserDataFromAPIGateWay struct {
	ClientId         string `json:"clientId" example:"5"`
	UserId           string `json:"userId" example:"5"`
	UserRole         string `json:"userRole" example:"admin"`
	Token            string `json:"token" example:"5"`
	RealmNameOfToken string `json:"realmNameOfToken" example:"master"`
	AuthCompleted    bool   `json:"authCompleted" example:"true"`
} //@name UserDataFromAPIGateWay

// {"clientId": "", "userId": "", "token":"", "authCompleted":true, "realmNameOfToken":"master"}

type IdRequest struct {
	ID int `json:"id" binding:"required" example:"1"`
} //@name IdRequest

type DeleteRequest struct {
	Ids []string `json:"ids" binding:"required" example:"[1,2,3]"`
} //@name DeleteRequest

type GetRequest struct {
	Skip  int `json:"skip" example:"0"`   // records to skip
	Limit int `json:"limit" example:"10"` // limit for the records. non of records to fetch per page
} //@name GetRequest

type Stamps struct {
	CreatedAt int64 `json:"createdAt" bson:"createdAt" example:"1600000000"`
	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt" example:"1600000000"`
} //@name Stamps

type Address struct {
	Id                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;index" json:"id" example:"a577055d-f40a-4617-9dc4-a6a81b317c8b"`
	Line1               string    `gorm:"column:line1" json:"line1" binding:"required" example:"A/604"`
	Line2               string    `gorm:"column:line2" json:"line2" example:"MaximaLines"`
	AreaOfDistrict      string    `gorm:"column:area_of_district" json:"areaOfDistrict" example:"Virar"`
	City                string    `gorm:"column:city" json:"city" binding:"required" example:"Mumbai"`
	State               string    `gorm:"column:state" json:"state" binding:"required" example:"Maharashtra"`
	Country             string    `gorm:"column:country" json:"country" binding:"required" example:"India"`
	ZipCode             string    `gorm:"column:zip_code" json:"zipCode" binding:"required" example:"401203"`
	Lat                 float64   `gorm:"column:lat" json:"lat" example:"19.423656"`
	Long                float64   `gorm:"column:long" json:"long" example:"72.819238"`
	GooglePlaceId       string    `gorm:"column:google_place_id" json:"googlePlaceId" example:"googlePlaceId"`
	HouseBuildingNumber string    `gorm:"column:house_building_number" json:"houseBuildingNumber" example:"604"`
	Email               string    `gorm:"column:email" json:"email" example:"xyz@gmail.com"`
	Phone               string    `gorm:"column:phone" json:"phone" example:"1234567890"`
	Website             string    `gorm:"column:website" json:"website" example:"www.xyz.com"`
	CountryCode         string    `gorm:"column:country_code" json:"countryCode" example:"+91"`
	AddressType         string    `gorm:"column:address_type" json:"addressType" example:"businessAddress"` // businessAddress, billingAddress, propertyAddress
} //@name Address

type Fields struct {
	Id          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;index" json:"fieldId" example:"a577055d-f40a-4617-9dc4-a6a81b317c8b"`
	FieldName   string    `gorm:"column:field_name" json:"fieldName" example:"IBAN"`
	Value       string    `gorm:"column:value" json:"value" example:"123456789"`
	FieldType   string    `gorm:"column:field_type" json:"fieldType" example:"string"` //bankAdditionalField, taxField
	IsMandatory bool      `gorm:"column:is_mandatory" json:"isMandatory" example:"true"`
	Description string    `gorm:"column:description" json:"description" example:"discription"`
} //@name Fields

type TokenResponse struct {
	ClientName          string                 `json:"ClientName"`
	CLientId            string                 `json:"CLientId"`
	UserId              uuid.UUID              `json:"UserId"`
	KeycloakUserId      string                 `json:"keycloakUserId"`
	UserType            string                 `json:"userType"`
	InstitutionCategory []string               `json:"institutionCategory"`
	TimeZone            string                 `json:"timeZone"`
	CompanyName         string                 `json:"companyName"`
	Name                string                 `json:"name"`
	PhoneNumber         string                 `json:"phoneNumber"`
	Email               string                 `json:"email"`
	Permission          map[string]interface{} `json:"permission"`
	Token               interface{}            `json:"token"`
} //@name TokenResponse

type StatusLogs struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id" example:"a577055d-f40a-4617-9dc4-a6a81b317c8b"`
	Status           Status    `gorm:"column:status;type:status_type" json:"status" example:"active"`
	CampaignID       uuid.UUID `gorm:"column:campaign_id" json:"campaignId" example:"5f5f5f5f5f5f5f5f5f5f5f5f"` // 5f5f5f5f5f5f5f5f5f5f5f5f
	ActionByUserRole string    `gorm:"column:action_by_user_role" json:"actionByUserRole" example:"admin"`      // admin, user
	ActionByUserId   string    `gorm:"column:action_by_user_id" json:"actionByUserId" example:"1"`              // 1, 2
	Notes            string    `gorm:"column:notes" json:"note" example:"xyz"`                                  // xyz
	Timestamp        int64     `gorm:"column:timestamp" json:"timestamp"`                                       // 2020-09-01T00:00:00Z
} //@name StatusLogs

type Logs struct {
	Id         uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id" example:"a577055d-f40a-4617-9dc4-a6a81b317c8b"`
	Trigger    string      `gorm:"column:trigger" json:"trigger" example:"create"`                                       // created, updated, deleted
	Entity     string      `gorm:"column:entity" json:"entity" example:"property"`                                       // property, user, client
	EntityId   string      `gorm:"column:entity_id" json:"entityId" example:"5f5f5f5f5f5f5f5f5f5f5f5f"`                  // 5f5f5f5f5f5f5f5f5f5f5f5f
	ClientName string      `gorm:"column:client_name" json:"clientName" example:"clientName"`                            // clientName
	ActionById string      `gorm:"column:action_by_id" json:"actionById" example:"a577055d-f40a-4617-9dc4-a6a81b317c8b"` // id of user who performed the action
	OldData    interface{} `gorm:"column:old_data;type:jsonb" json:"oldData" example:"oldData"`                          // oldData
	NewData    interface{} `gorm:"column:new_data;type:jsonb" json:"newData" example:"newData"`                          // newData
	Timestamp  time.Time   `gorm:"column:timestamp" json:"timestamp" example:"1600000000"`                               // 1600000000
} //@name Logs

type Media struct {
	Id         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id" example:"a577055d-f40a-4617-9dc4-a6a81b317c8b"`
	MediaType  string    `gorm:"column:media_type" json:"mediaType" example:"image"`
	Url        string    `gorm:"column:url" json:"url" example:"https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"`
	AltText    string    `gorm:"column:alt_text" json:"altText" example:"google logo"`
	Position   int       `gorm:"column:position" json:"position" example:"1"`
	PropertyId string    `gorm:"column:property_id" json:"propertyId" example:"1"`
	ProductId  string    `gorm:"column:product_id" json:"productId" example:"1"`
} //@name Media
