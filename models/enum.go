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

type CampaignStatus int

// Campaign status constants
const (
	CampaignStatusDraft     = 0
	CampaignStatusActive    = 1
	CampaignStatusCompleted = 2
	CampaignStatusCancelled = 3
	CampaignStatusFull      = 4
)

type ParticipantStatus int

// Participant status constants
const (
	ParticipantStatusPending  = 0
	ParticipantStatusActive   = 1
	ParticipantStatusLeft     = 2
	ParticipantStatusRejected = 3
)
