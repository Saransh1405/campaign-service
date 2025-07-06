package models

//login

type LoginRequest struct {
	ClientName  string `json:"clientName" example:"zee"`
	Email       string `json:"email" example:"xyz@gmail.com"`
	Password    string `json:"password" example:"Xyz@Company"`
	Phone       string `json:"phone" example:"1234567890"`
	CountryCode string `json:"countryCode" example:"+91"`
	Otp         string `json:"otp" example:"123456"`
} //@name LoginRequest

//change password

type PatchPasswordRequest struct {
	Email              string `json:"email" example:"xyz@gmail.com"`
	OldPassword        string `json:"oldPassword" binding:"required" example:"Xyz@Company"`
	NewPassword        string `json:"newPassword" binding:"required" example:"Xyz@CORPO"`
	ConfirmNewPassword string `json:"confirmNewPassword" binding:"required" example:"Xyz@CORPO"`
} //@name PatchPasswordRequest

//reset password

type PostResetPasswordRequest struct {
	Email              string `json:"email" bson:"email" binding:"required" example:"xyz@gmail.com"`
	NewPassword        string `json:"newPassword" bson:"newPassword" binding:"required" example:"Xyz@CORPO"`
	ConfirmNewPassword string `json:"confirmNewPassword" bson:"confirmNewPassword" binding:"required" example:"Xyz@CORPO"`
} //@name PostResetPasswordRequest
