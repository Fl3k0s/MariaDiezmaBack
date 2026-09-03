package domain

import (
	"time"
)

type RequestStatus string

const (
	StatusPending    RequestStatus = "pending"
	StatusInProgress RequestStatus = "in_progress"
	StatusResolved   RequestStatus = "resolved"
	StatusCancelled  RequestStatus = "cancelled"
)

type RequestPriority string

const (
	PriorityLow    RequestPriority = "low"
	PriorityMedium RequestPriority = "medium"
	PriorityHigh   RequestPriority = "high"
	PriorityUrgent RequestPriority = "urgent"
)

type RequestItem struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`                  // e.g. "contact", "budget", "booking", "inquiry"
	Status      RequestStatus   `json:"status"`                // pending, in_progress, resolved, cancelled
	Priority    RequestPriority `json:"priority"`              // low, medium, high, urgent
	SenderName  string          `json:"sender_name"`
	SenderEmail string          `json:"sender_email"`
	SenderPhone string          `json:"sender_phone,omitempty"`
	Subject     string          `json:"subject"`
	Message     string          `json:"message"`
	Metadata    map[string]any  `json:"metadata,omitempty"`
	InternalNote string         `json:"internal_note,omitempty"`
	AssignedTo  *string         `json:"assigned_to,omitempty"` // Admin/Operator user ID
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateRequestInput struct {
	Type        string          `json:"type"`
	Priority    RequestPriority `json:"priority"`
	SenderName  string          `json:"sender_name"`
	SenderEmail string          `json:"sender_email"`
	SenderPhone string          `json:"sender_phone,omitempty"`
	Subject     string          `json:"subject"`
	Message     string          `json:"message"`
	Metadata    map[string]any  `json:"metadata,omitempty"`
}

type UpdateRequestInput struct {
	Type         *string          `json:"type,omitempty"`
	Status       *RequestStatus   `json:"status,omitempty"`
	Priority     *RequestPriority `json:"priority,omitempty"`
	Subject      *string          `json:"subject,omitempty"`
	Message      *string          `json:"message,omitempty"`
	InternalNote *string          `json:"internal_note,omitempty"`
	AssignedTo   *string          `json:"assigned_to,omitempty"`
}

type UpdateStatusInput struct {
	Status       RequestStatus `json:"status"`
	InternalNote string        `json:"internal_note,omitempty"`
}

type RequestFilter struct {
	Type      string
	Status    RequestStatus
	Priority  RequestPriority
	Search    string // search across sender_name, sender_email, subject
	Page      int
	PerPage   int
	SortBy    string // "created_at", "priority", "status"
	SortOrder string // "asc", "desc"
}
