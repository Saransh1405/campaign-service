package models

type JoinCampaignRequest struct {
	CampaignID string `json:"campaign_id" binding:"required"`
	UserID     string `json:"user_id" binding:"required"`
}

type LeaveCampaignRequest struct {
	CampaignID string `json:"campaign_id" binding:"required"`
	UserID     string `json:"user_id" binding:"required"`
}
