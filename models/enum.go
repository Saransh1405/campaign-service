package models

type Status string

const (
	Active          Status = "Active"
	Suspended       Status = "Suspended"
	Deleted         Status = "Deleted"
	Inactive        Status = "Inactive"
	PendindApproval Status = "Pending Approval"
	Rejected        Status = "Rejected"
	Approved        Status = "Approved"
	Submitted       Status = "Submitted"
)

type CampaignStatus string

// Campaign status constants
const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusActive    CampaignStatus = "active"
	CampaignStatusCompleted CampaignStatus = "completed"
	CampaignStatusCancelled CampaignStatus = "cancelled"
	CampaignStatusFull      CampaignStatus = "full"
)

type ParticipantStatus string

// Participant status constants
const (
	ParticipantStatusPending  ParticipantStatus = "pending"
	ParticipantStatusActive   ParticipantStatus = "active"
	ParticipantStatusLeft     ParticipantStatus = "left"
	ParticipantStatusRejected ParticipantStatus = "rejected"
)

type CampaignEventType string

const (
	CreateCampaignEvent CampaignEventType = "create_campaign"
	UpdateCampaignEvent CampaignEventType = "update_campaign"
	DeleteCampaignEvent CampaignEventType = "delete_campaign"
)

type CampaignActivityType string

const (
	CampaignActivity CampaignActivityType = "campaign_activity"
)
