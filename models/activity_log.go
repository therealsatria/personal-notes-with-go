package models

import "time"

// ActivityLog represents an activity log entry in the system
type ActivityLog struct {
	ID          int       `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Action      string    `json:"action"`     // create, read, update, delete, check, generate
	EntityType  string    `json:"entityType"` // note, category, encryption, key
	EntityID    int       `json:"entityId"`
	Description string    `json:"description"`
	UserID      int       `json:"userId"`
	IPAddress   string    `json:"ipAddress"`
}

// ActivityLogFilter represents filters for activity logs
type ActivityLogFilter struct {
	EntityType string
	Action     string
	StartDate  time.Time
	EndDate    time.Time
	Limit      int
	Offset     int
}
