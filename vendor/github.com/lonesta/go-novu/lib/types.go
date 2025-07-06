package lib

type ChannelType string
type TypeEnum string

// TypeData type: 'data' => will turn this into a FCM data notification, where the payload is sent as a data notification
// TypeNotification If type is 'notification', you can use the "data" (Fcm) override to send notification messages with optional data payload
const (
	TypeData         TypeEnum = "data"
	TypeNotification TypeEnum = "notification"
)

func (e TypeEnum) String() string {
	return string(e)
}

type Fcm struct {
	Type TypeEnum    `json:"type"`
	Data interface{} `json:"data"`
}

type Overrides struct {
	Fcm Fcm `json:"fcm"`
}

type ChannelCredentials struct {
	WebhookUrl   string   `json:"webhookUrl"`
	DeviceTokens []string `json:"deviceTokens"`
}
