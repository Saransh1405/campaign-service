package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Campaign struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      string    `json:"user_id" gorm:"type:string;not null"`
	Name        string    `json:"name" gorm:"not null;size:255"`
	Description string    `json:"description" gorm:"type:text"`
	Type        string    `json:"type" gorm:"not null;size:100"` // cricket, football, etc.
	ImageURL    string    `json:"image_url" gorm:"size:500"`
	DisplayName string    `json:"display_name" gorm:"size:255"`

	// Timing
	StartDate int64 `json:"start_date" gorm:"not null"`
	EndDate   int64 `json:"end_date" gorm:"not null"`

	// Capacity
	MaxParticipants int `json:"max_participants" gorm:"not null"`
	MinParticipants int `json:"min_participants" gorm:"not null"`
	CurrentCount    int `json:"current_count" gorm:"default:0"`

	// Pricing
	Price    int64  `json:"price" gorm:"not null"`
	Currency string `json:"currency" gorm:"default:'INR';size:3"`

	// Status and visibility
	Status   CampaignStatus `json:"status" gorm:"type:campaign_status;default:'draft'"`
	IsPublic bool           `json:"is_public" gorm:"default:true"`

	// Metadata
	Tags      json.RawMessage `json:"tags" gorm:"type:jsonb"`
	CreatedAt int64           `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64           `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	StatusLogs   []StatusLogs  `json:"status_logs" gorm:"foreignKey:CampaignID"`
	Categories   []Category    `json:"categories" gorm:"many2many:campaign_categories"`
	Participants []Participant `json:"participants" gorm:"foreignKey:CampaignID"`
	Location     *Location     `json:"location" gorm:"foreignKey:CampaignID"`
}

type Location struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CampaignID uuid.UUID `json:"campaign_id" gorm:"type:uuid;not null"`
	Name       string    `json:"name" gorm:"size:255"` // Venue name
	Address    string    `json:"address" gorm:"size:500"`
	Latitude   float64   `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude  float64   `json:"longitude" gorm:"type:decimal(11,8)"`
	City       string    `json:"city" gorm:"size:100"`
	State      string    `json:"state" gorm:"size:100"`
	Country    string    `json:"country" gorm:"size:100"`
	ZipCode    string    `json:"zip_code" gorm:"size:20"`
	CreatedAt  int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  int64     `json:"updated_at" gorm:"autoUpdateTime"`
}

type Participant struct {
	ID         uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID         `json:"user_id" gorm:"type:uuid;not null"`
	CampaignID uuid.UUID         `json:"campaign_id" gorm:"type:uuid;not null"`
	JoinedAt   int64             `json:"joined_at" gorm:"autoCreateTime"`
	LeftAt     int64             `json:"left_at,omitempty"`
	Status     ParticipantStatus `json:"status" gorm:"type:participant_status;default:'pending'"`
	PaymentID  *string           `json:"payment_id,omitempty" gorm:"size:100"` // Reference to payment
	CreatedAt  int64             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  int64             `json:"updated_at" gorm:"autoUpdateTime"`
}

type Category struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:100"`
	DisplayName string    `json:"display_name" gorm:"size:100"`
	ImageURL    string    `json:"image_url" gorm:"size:500"`
	Description string    `json:"description" gorm:"type:text"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at" gorm:"autoUpdateTime"`
}

// Additional models for advanced features
type CampaignInvite struct {
	ID         uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CampaignID uuid.UUID         `json:"campaign_id" gorm:"type:uuid;not null"`
	InvitedBy  uuid.UUID         `json:"invited_by" gorm:"type:uuid;not null"`
	InvitedTo  uuid.UUID         `json:"invited_to" gorm:"type:uuid;not null"`
	Status     ParticipantStatus `json:"status" gorm:"type:participant_status;default:'pending'"` // pending, accepted, rejected
	ExpiresAt  int64             `json:"expires_at"`
	CreatedAt  int64             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  int64             `json:"updated_at" gorm:"autoUpdateTime"`
}

type CampaignReview struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CampaignID uuid.UUID `json:"campaign_id" gorm:"type:uuid;not null"`
	UserID     string    `json:"user_id" gorm:"type:uuid;not null"`
	Rating     int       `json:"rating" gorm:"check:rating >= 1 AND rating <= 5"`
	Comment    string    `json:"comment" gorm:"type:text"`
	CreatedAt  int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  int64     `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateCampaignRequest struct {
	UserID          string          `json:"user_id" binding:"required"`
	Name            string          `json:"name" binding:"required,min=3,max=255"`
	Description     string          `json:"description" binding:"required,min=10"`
	Type            string          `json:"type" binding:"required"`
	ImageURL        string          `json:"image_url" binding:"url"`
	DisplayName     string          `json:"display_name"`
	StartDate       int64           `json:"start_date" binding:"required"`
	EndDate         int64           `json:"end_date" binding:"required"`
	MaxParticipants int             `json:"max_participants" binding:"required,min=2"`
	MinParticipants int             `json:"min_participants" binding:"required,min=1"`
	Price           int64           `json:"price" binding:"min=0"`
	Currency        string          `json:"currency" binding:"required,len=3"`
	IsPublic        bool            `json:"is_public"`
	Tags            json.RawMessage `json:"tags"`
	Location        Location        `json:"location" binding:"required"`
}

type UpdateCampaignRequest struct {
	ID              *string         `json:"id" binding:"required"`
	UserID          string          `json:"user_id" binding:"required"`
	Name            *string         `json:"name,omitempty"`
	Description     *string         `json:"description,omitempty"`
	ImageURL        *string         `json:"image_url,omitempty"`
	DisplayName     *string         `json:"display_name,omitempty"`
	StartDate       *int64          `json:"start_date,omitempty"`
	EndDate         *int64          `json:"end_date,omitempty"`
	MaxParticipants *int            `json:"max_participants,omitempty"`
	MinParticipants *int            `json:"min_participants,omitempty"`
	Price           *int64          `json:"price,omitempty"`
	IsPublic        *bool           `json:"is_public,omitempty"`
	Tags            json.RawMessage `json:"tags,omitempty"`
	Location        *Location       `json:"location,omitempty"`
	Status          *CampaignStatus `json:"status,omitempty"`
}

type GetCampaignRequest struct {
	GetRequest
	ID        string   `form:"id"`
	Status    string   `form:"status"`
	UserID    string   `form:"user_id"`
	City      string   `form:"city"`
	State     string   `form:"state"`
	Country   string   `form:"country"`
	MinPrice  int64    `form:"min_price"`
	MaxPrice  int64    `form:"max_price"`
	StartDate string   `form:"start_date"` // YYYY-MM-DD
	EndDate   string   `form:"end_date"`   // YYYY-MM-DD
	Tags      []string `form:"tags"`
	SortBy    string   `form:"sort_by"`    // start_date, created_at, price, participants
	SortOrder string   `form:"sort_order"` // asc, desc
}

type CampaignEvent struct {
	Campaign         map[string]interface{} `json:"campaign"`
	UpdateFields     map[string]interface{} `json:"update_fields"`
	EventPublishTime int64                  `json:"event_publish_time"`
	EventType        string                 `json:"event_type"`
}
