package models

// verify email
type PostVerifyEmailRequest struct {
	Email string `json:"email" bson:"email" binding:"required" example:"xyz@gmail.com"`
} //@name PostVerifyEmailRequest

// send verify email
type PostSendVerifyEmailRequest struct {
	Email string `json:"email" bson:"email" binding:"required" example:"xyz@gmail.com"`
} //@name PostSendVerifyEmailRequest

// change email
type PatchVerifyEmailRequest struct {
	Otp string `json:"otp" bson:"otp" binding:"required" example:"436457"`
} //@name PatchVerifyEmailRequest

type VerifyPhoneRequest struct {
	Phone       string `json:"phone" binding:"required" example:"1234567890"`
	CountryCode string `json:"countryCode" binding:"required" example:"+91"`
	ClientName  string `json:"clientName" binding:"required" example:"zee5"`
	Otp         string `json:"otp" binding:"required" example:"436457"`
} //@name VerifyPhoneRequest
