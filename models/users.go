package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"5f5f5f5f5f5f5f5f5f5f5f5f"`
	SuspendedAt         int64              `bson:"suspendedAt"  json:"suspendedAt"`
	FirstName           string             `bson:"firstName" example:"John" json:"firstName"`
	LastName            string             `bson:"lastName" example:"Doe" json:"lastName"`
	Password            string             `bson:"password" example:"password123" json:"password"`
	Email               string             `bson:"email" example:"jhon@example.com" json:"email"`
	CountryCode         string             `bson:"countryCode" example:"+91" json:"countryCode"`
	Phone               string             `bson:"phone" example:"1234567890" json:"phone"`
	PhoneVerified       bool               `bson:"phoneVerified" json:"phoneVerified"`
	UserProfileUrl      string             `bson:"userProfileUrl" example:"www.xyz.com" json:"userProfileUrl"`
	TemporaryPassword   string             `bson:"temporaryPassword" example:"123456" json:"temporaryPassword"`
	Status              Status             `bson:"status"  example:"Active" json:"status"`
	ReasonForDeletion   string             `bson:"reasonForDeletion" example:"not verified" json:"reasonForDeletion"`
	ReasonForSuspension string             `bson:"reasonForSuspension" example:"not verified" json:"reasonForSuspension"`
	StatusLogs          []StatusLogs       `bson:"statusLogs" json:"statusLogs"`
	ClientName          string             `bson:"clientName" json:"clientName"`
	CreatedAt           int64              `bson:"createdAt" json:"createdAt"`
	// Google OAuth fields
	GoogleID      string `bson:"googleId,omitempty" json:"googleId,omitempty" example:"123456789"`
	AuthProvider  string `bson:"authProvider,omitempty" json:"authProvider,omitempty" example:"google"`
	EmailVerified bool   `bson:"emailVerified,omitempty" json:"emailVerified,omitempty" example:"true"`
} //@name Users
