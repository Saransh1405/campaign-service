package constants

// Error Codes
const (
	RequestBodyBindError        = "request body bind error"
	RequestBodyValidationError  = "request body validation error"
	ExternalServiceFailureError = "external service failure error"
	DatabaseFailureError        = "database failure error"
	RequestValidationError      = "request validation error"
	BadRequestError             = "bad request error"
	UnauthorizedMessage         = "common.401"
	SuccessMessage              = "common.200"
	InternalServerMessage       = "common.500"
	BadRequestMessage           = "common.400"
	NotFoundMessage             = "common.404"
	ConflictMessage             = "common.409"
	AccountDisableMessage       = "common.422"
	RetryWithMessage            = "common.401"
	TooManyRequestsMessage      = "common.429"
	NotFound                    = "no data found"
	AlreadyExists               = "data Already Exists"
	PasswordUpdated             = "password updated successfully"
	SuccessMsg                  = "success"

	DirectAccessNotAllowed = "direct access not allowed"
	BadRequestErr          = "the server could not understand the request that it was sent."
)

const (
	DataDeletedSuccess = "Data Deleted Successfully"
)

// ErrorCode

const (
	ErrorCode = "1"
)

const (
	// WLBadRequestCode is the code for bad request
	WLBadRequestCode = "400"

	// WLUnauthorizedCode is the code for unauthorized
	WLUnauthorizedCode = "401"

	// WLForbiddenCode is the code for payment required
	WLPaymentRequiredCode = "402"

	// WLForbiddenCode is the code for user account not active
	WLForbiddenCode = "403"

	// WLRateLimitCode is the code for rate limit
	WLRateLimitCode = "429"

	// WLConflictCode is the code for conflict
	WLDataConflictCode = "409"

	// WLNoDataFoundCode is the code for no data found
	WLNoDataFoundCode = "404"

	// WLInternalServerErrorCode is the code for internal server error
	WLInternalServerErrorCode = "500"

	// WLAccountDisableCode is the code for account disable
	WLAccountDisableCode = "403"

	//WLSuccessCode is the code for success
	WLSuccessCode = "200"
)

// Error Messages
const (
	DataNotDeletedMessage                = "errors.DataNotDeleted"
	DataNotUpdatedMessage                = "errors.DataNotUpdated"
	AccountsDataNotFoundMessage          = "errors.AccountsDataNotFound"
	RoleNotFoundMessage                  = "errors.RoleNotFound"
	UserNotFoundMessage                  = "errors.UserNotFound"
	UserNotActiveMessage                 = "errors.UserNotActive"
	UserAlreadyExistsMessage             = "errors.UserAlreadyExists"
	SupportedLanguageNotFoundMessage     = "errors.SupportedLanguageNotFound"
	CurrencyNotFoundMessage              = "errors.CurrencyNotFound"
	StateTaxNotFoundMessage              = "errors.StateTaxNotFound"
	AccountTypeNotFoundMessage           = "errors.AccountTypeNotFound"
	InvalidOldPasswordMessage            = "errors.InvalidOldPassword"
	PasswordDoesNotMatchMessage          = "errors.PasswordDoesNotMatch"
	OldPasswordAndNewPasswordSameMessage = "errors.OldPasswordAndNewPasswordSame"
	BusinessIdIsRequiredMessage          = "errors.BusinessIdIsRequired"
	BrandIdIsRequiredMessage             = "errors.BrandIdIsRequired"
	InvalidOTPMessage                    = "errors.InvalidOTP"
)

// Trigger Messages
const (
	CreatedTrigger = "Created"
	UpdatedTrigger = "Updated"
	DeletedTrigger = "Deleted"
)
