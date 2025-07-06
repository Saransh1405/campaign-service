package constants

// common constants
const (
	ApplicationName = "users-service"
)

const (
	EnumKey = "enum"
)

const (
	ABCDLower                     = "abcdefghijklmnopqrstuvwxyz"
	ABCDUpper                     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Number                        = "0123456789"
	SpecialCharSet2               = "!@#$%&*"
	RequestId                     = "requestId"
	StatusCode                    = "statusCode"
	StartTimeLogParam             = "startTime"
	QueryLogParam                 = "query"
	HeaderLogParams               = "header"
	LatencyLogParam               = "latency"
	ClientIPLogParam              = "clientIP"
	MethodLogParam                = "method"
	TimeLogParam                  = "time"
	ErrorLogParam                 = "error"
	RequestIDHeader               = "X-requestId"
	IST                           = "Asia/Kolkata"
	IsJsonArray                   = 0
	IsString                      = 1
	APIRespErrorKey               = "Request not completed"
	APIRespSuccessKey             = "Request successfully completed"
	APIRespConflictKey            = "gg"
	ChartRequestPayload           = "hh"
	ErrConflictMSG                = "Conflict"
	WLConflictCode                = "1000"
	TimeFormat                    = "time.UnixDate"
	BindingFailedErrr             = "Binding fail"
	UnauthorizedErrr              = "Unauthorized"
	NotFoundErrr                  = "Not Found"
	VaultInitializationFailedErrr = "Vault initialization failed"
	LanguageString                = "language"
	LogoutSuccessMessage          = "Logout successfully"
)

const (
	BillingAddress      = "Billing Address"
	BusinessAddress     = "Business Address"
	PropertyAddress     = "Property Address"
	BankAdditionalField = "Bank Additional Field"
	TaxField            = "Tax Field"
)

// user types

const (
	BusinessAdminType = "Business Admin"
	BrandAdminType    = "Brand Admin"
	PropertyAdminType = "Property Admin"
)
