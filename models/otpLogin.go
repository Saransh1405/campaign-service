package models

type SendOTPRequest struct {
	ClientName  string `json:"clientName"`
	CountryCode string `json:"countryCode" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
} // @name SendOTPRequest

type LoginWithOTPRequest struct {
	ClientName  string `json:"clientName" binding:"required"`
	CountryCode string `json:"countryCode" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
} // @name LoginWithOTPRequest

type LoginResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	SessionState     string `json:"session_state"`
	TokenType        string `json:"token_type"`
} // @name LoginResponse
