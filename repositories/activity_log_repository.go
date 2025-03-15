package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"personal-notes-with-go/models"
)

// ActivityLogRepository handles database operations for activity logs
type ActivityLogRepository struct {
	DB *sql.DB
}

// NewActivityLogRepository creates a new ActivityLogRepository
func NewActivityLogRepository(db *sql.DB) *ActivityLogRepository {
	return &ActivityLogRepository{DB: db}
}

// CreateTable creates the activity_logs table if it doesn't exist
func (r *ActivityLogRepository) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS activity_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME NOT NULL,
		action TEXT NOT NULL,
		entity_type TEXT NOT NULL,
		entity_id INTEGER NOT NULL,
		description TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		ip_address TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_activity_logs_timestamp ON activity_logs(timestamp);
	CREATE INDEX IF NOT EXISTS idx_activity_logs_entity_type ON activity_logs(entity_type);
	CREATE INDEX IF NOT EXISTS idx_activity_logs_action ON activity_logs(action);
	`

	_, err := r.DB.Exec(query)
	return err
}

// Create adds a new activity log entry
func (r *ActivityLogRepository) Create(log *models.ActivityLog) error {
	query := `
	INSERT INTO activity_logs (
		timestamp, action, entity_type, entity_id, description, user_id, ip_address
	) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	result, err := r.DB.Exec(
		query,
		log.Timestamp,
		log.Action,
		log.EntityType,
		log.EntityID,
		log.Description,
		log.UserID,
		log.IPAddress,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	log.ID = int(id)
	return nil
}

// GetAll retrieves all activity logs with optional filtering
func (r *ActivityLogRepository) GetAll(filter models.ActivityLogFilter) ([]models.ActivityLog, error) {
	query := `
	SELECT id, timestamp, action, entity_type, entity_id, description, user_id, ip_address
	FROM activity_logs
	WHERE 1=1
	`
	var args []interface{}

	if filter.EntityType != "" {
		query += " AND entity_type = ?"
		args = append(args, filter.EntityType)
	}

	if filter.Action != "" {
		query += " AND action = ?"
		args = append(args, filter.Action)
	}

	if !filter.StartDate.IsZero() {
		query += " AND timestamp >= ?"
		args = append(args, filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		query += " AND timestamp <= ?"
		args = append(args, filter.EndDate)
	}

	query += " ORDER BY timestamp DESC"

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)

		if filter.Offset > 0 {
			query += " OFFSET ?"
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ActivityLog
	for rows.Next() {
		var log models.ActivityLog
		err := rows.Scan(
			&log.ID,
			&log.Timestamp,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&log.Description,
			&log.UserID,
			&log.IPAddress,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetByID retrieves an activity log by ID
func (r *ActivityLogRepository) GetByID(id int) (*models.ActivityLog, error) {
	query := `
	SELECT id, timestamp, action, entity_type, entity_id, description, user_id, ip_address
	FROM activity_logs
	WHERE id = ?
	`

	var log models.ActivityLog
	err := r.DB.QueryRow(query, id).Scan(
		&log.ID,
		&log.Timestamp,
		&log.Action,
		&log.EntityType,
		&log.EntityID,
		&log.Description,
		&log.UserID,
		&log.IPAddress,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("activity log with ID %d not found", id)
		}
		return nil, err
	}

	return &log, nil
}

// DeleteOlderThan deletes activity logs older than the specified days
func (r *ActivityLogRepository) DeleteOlderThan(days int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	query := "DELETE FROM activity_logs WHERE timestamp < ?"

	result, err := r.DB.Exec(query, cutoffDate)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// LogActivity is a helper function to create an activity log entry
func (r *ActivityLogRepository) LogActivity(action, entityType string, entityID int, description string, userID int, ipAddress string) {
	activityLog := &models.ActivityLog{
		Timestamp:   time.Now(),
		Action:      action,
		EntityType:  entityType,
		EntityID:    entityID,
		Description: description,
		UserID:      userID,
		IPAddress:   ipAddress,
	}

	err := r.Create(activityLog)
	if err != nil {
		log.Printf("Failed to log activity: %v", err)
	}
}

// Count returns the total number of activity logs with optional filtering
func (r *ActivityLogRepository) Count(filter models.ActivityLogFilter) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM activity_logs
	WHERE 1=1
	`
	var args []interface{}

	if filter.EntityType != "" {
		query += " AND entity_type = ?"
		args = append(args, filter.EntityType)
	}

	if filter.Action != "" {
		query += " AND action = ?"
		args = append(args, filter.Action)
	}

	if !filter.StartDate.IsZero() {
		query += " AND timestamp >= ?"
		args = append(args, filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		query += " AND timestamp <= ?"
		args = append(args, filter.EndDate)
	}

	var count int
	err := r.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
