package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPostRequest struct {
	FirstName         string   `json:"firstName" binding:"required" example:"John"`
	LastName          string   `json:"lastName" binding:"required" example:"Doe"`
	Email             string   `json:"email" binding:"required" example:"xyz@gmail.com"`
	Phone             string   `json:"phone" binding:"required" example:"1234567890"`
	CountryCode       string   `json:"countryCode" binding:"required" example:"+91"`
	ProfilePictureUrl string   `json:"profilePictureUrl" example:"www.xyz.com"`
	UserType          string   `json:"userType" example:"admin"`
	RoleId            []string `json:"roleId" example:"admin"`
} //@name UserPostRequest

type UserPatchRequest struct {
	Id                  string `json:"id" binding:"required" example:"5f5f5f5f5f5f5f5f5f5f5f5f"`
	FirstName           string `json:"firstName" example:"John"`
	LastName            string `json:"lastName" example:"Doe"`
	ProfilePictureUrl   string `json:"profilePictureUrl" example:"www.xyz.com"`
	Status              string `json:"status" example:"active"`
	ReasonForSuspension string `json:"reasonForSuspension" example:"xyz"`
} //@name UserPatchRequest

type UserDeleteRequest struct {
	ID                []string `json:"id" binding:"required" example:"5"`
	ReasonForDeletion string   `json:"reasonForDeletion" bson:"reasonForDeletion" example:"not verified"`
} //@name UserDeleteRequest

type Users struct {
	Model
	SuspendedAt         *time.Time   `gorm:"column:suspended_at" example:"2020-09-01T00:00:00Z" json:"suspendedAt"`
	FirstName           string       `gorm:"column:first_name" example:"John" json:"firstName"`
	LastName            string       `gorm:"column:last_name" example:"Doe" json:"lastName"`
	Email               string       `gorm:"column:email;index:email_client_name;index:email_keycloak_user_id_client_name_phone_country_code;uniqueIndex:email_keycloak_user_id_client_name" example:"jhon@example.com" json:"email"`
	CountryCode         string       `gorm:"column:country_code;index:email_keycloak_user_id_client_name_phone_country_code" example:"+91" json:"countryCode"`
	Phone               string       `gorm:"column:phone;index:email_keycloak_user_id_client_name_phone_country_code;uniqueIndex:email_keycloak_user_id_client_name;uniqueIndex:phone_client_name" example:"1234567890" json:"phone"`
	PhoneVerified       bool         `gorm:"column:phone_verified" json:"phoneVerified"`
	KeycloakUserId      string       `gorm:"column:keycloak_user_id;index:email_keycloak_user_id_client_name_phone_country_code;uniqueIndex:email_keycloak_user_id_client_name" example:"5f5f5f5f5f5f5f5f5f5f5f5f" json:"keycloakUserId"`
	ProfilePictureUrl   string       `gorm:"column:profile_picture_url" example:"www.xyz.com" json:"profilePictureUrl"`
	TemporaryPassword   string       `gorm:"column:temporary_password" example:"123456" json:"temporaryPassword"`
	Status              Status       `gorm:"column:status; type:status" example:"Active" json:"status"`
	ReasonForDeletion   string       `gorm:"column:reason_for_deletion" example:"not verified" json:"reasonForDeletion"`
	ReasonForSuspension string       `gorm:"column:reason_for_suspension" example:"not verified" json:"reasonForSuspension"`
	StatusLogs          []StatusLogs `gorm:"foreignKey:EntityId;constraint:OnDelete:CASCADE"`
	UserType            string       `gorm:"column:user_type;index:client_name_user_type" example:"admin" json:"userType"`
	ClientName          string       `gorm:"column:client_name;index:email_keycloak_user_id_client_name_phone_country_code;uniqueIndex:email_keycloak_user_id_client_name;index:client_name_user_type;index;index:email_client_name;uniqueIndex:phone_client_name" example:"client1" json:"clientName"`
} //@name Users

type UserResponse struct {
	ID                  uuid.UUID  `json:"id" example:"5f5f5f5f5f5f5f5f5f5f5f5f"`                                          //User id
	CreatedAt           time.Time  `json:"createdAt" example:"2020-09-01T00:00:00Z"`                                       //User created at
	UpdatedAt           time.Time  `json:"updatedAt" example:"2020-09-01T00:00:00Z"`                                       //User updated at
	DeletedAt           time.Time  `json:"deletedAt" example:"2020-09-01T00:00:00Z"`                                       //User deleted at
	SuspendedAt         *time.Time `json:"suspendedAt" example:"2020-09-01T00:00:00Z"`                                     //User suspended at
	FirstName           string     `json:"firstName" example:"John"`                                                       //User first name
	LastName            string     `json:"lastName" example:"Doe"`                                                         //User last name
	Email               string     `json:"email" example:"jhon@gmail.com"`                                                 //User email
	CountryCode         string     `gorm:"column:country_code" example:"+91" json:"countryCode"`                           //User country code of phone number
	Phone               string     `json:"phone" example:"1234567890"`                                                     //User phone number
	PhoneVerified       bool       `json:"phoneVerified" example:"true"`                                                   //Check if user phone number is verified or not
	Status              Status     `json:"status" example:"active"`                                                        //User status active, deleted, suspended
	KeycloakUserId      string     `json:"keycloakUserId" example:"5f5f5f5f5f5f5f5f5f5f5f"`                                //User keycloak id
	ProfilePictureUrl   string     `json:"profilePictureUrl" example:"www.xyz.com"`                                        //User profile picture url
	UserType            string     `json:"userType" example:"admin"`                                                       //User type admin, employee, guest
	ClientName          string     `json:"clientName" example:"client1"`                                                   //Client name under which user is created
	ReasonForDeletion   string     `gorm:"column:reason_for_deletion" example:"not verified" json:"reasonForDeletion"`     //Reason for deletion of user
	ReasonForSuspension string     `gorm:"column:reason_for_suspension" example:"not verified" json:"reasonForSuspension"` //Reason for suspension of user
} //@name UserResponse
