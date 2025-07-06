package constants

type ApiKeyStatus string

const (
	Active   ApiKeyStatus = "Active"
	Inactive ApiKeyStatus = "Inactive"
	Delete   ApiKeyStatus = "Delete"
)

// IsValid Each enum must have this function for validation
func (c ApiKeyStatus) IsValid() bool {
	switch c {
	case Active, Inactive, Delete:
		return true
	}

	return false
}
