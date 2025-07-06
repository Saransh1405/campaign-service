package models

type SignupUserRequest struct {
	FirstName   string `json:"firstName" binding:"required" example:"John"`
	LastName    string `json:"lastName" binding:"required" example:"Doe"`
	CountryCode string `json:"countryCode" binding:"required" example:"91"`
	PhoneNumber string `json:"phoneNumber" binding:"required" example:"1234567890"`
} //@name SignupUserRequest
