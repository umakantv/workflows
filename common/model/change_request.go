package model

// ChangeRequestWorkflow is a Temporal workflow that processes a change request.
type ChangeRequest struct {
	ID         string              `gorm:"primaryKey" go:"primaryKey"`
	CreatedBy  string              `gorm:"column:created_by"`
	WorkflowID string              `gorm:"column:workflow_id"`
	RunID      string              `gorm:"column:run_id"`
	Status     ChangeRequestStatus `gorm:"column:status"`
	Approvers  string              `gorm:"column:approvers"`
	// Change     EntityChange        `gorm:"-"`
	// Existing   interface{}         `gorm:"-"`
}

// EntityChange represents the changes made to an entity.
type EntityChange struct {
	NewValue      interface{}
	PreviousValue interface{}
}

type ChangeRequestStatus string

const (
	ChangeRequestStatusDraft     ChangeRequestStatus = "DRAFT"
	ChangeRequestStatusDiscarded ChangeRequestStatus = "DISCARDED"
	ChangeRequestStatusSubmitted ChangeRequestStatus = "SUBMITTED"
	ChangeRequestStatusApproved  ChangeRequestStatus = "APPROVED"
	ChangeRequestStatusRejected  ChangeRequestStatus = "REJECTED"
)
