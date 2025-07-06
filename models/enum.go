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
