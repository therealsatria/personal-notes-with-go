package handlers

import (
	"net/http"
	"personal-notes-with-go/models"
	"personal-notes-with-go/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ActivityLogHandler handles HTTP requests for activity logs
type ActivityLogHandler struct {
	repo *repositories.ActivityLogRepository
}

// NewActivityLogHandler creates a new ActivityLogHandler
func NewActivityLogHandler(repo *repositories.ActivityLogRepository) *ActivityLogHandler {
	return &ActivityLogHandler{repo: repo}
}

// GetLogs handles GET /activity-logs
func (h *ActivityLogHandler) GetLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filter := models.ActivityLogFilter{
		Limit:  limit,
		Offset: offset,
	}

	logs, err := h.repo.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetLogsByEntityType handles GET /activity-logs/entity-type/:entityType
func (h *ActivityLogHandler) GetLogsByEntityType(c *gin.Context) {
	entityType := c.Param("entityType")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filter := models.ActivityLogFilter{
		EntityType: entityType,
		Limit:      limit,
		Offset:     offset,
	}

	logs, err := h.repo.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetLogsByAction handles GET /activity-logs/action/:action
func (h *ActivityLogHandler) GetLogsByAction(c *gin.Context) {
	action := c.Param("action")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filter := models.ActivityLogFilter{
		Action: action,
		Limit:  limit,
		Offset: offset,
	}

	logs, err := h.repo.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// DeleteOldLogs handles DELETE /activity-logs/older-than/:days
func (h *ActivityLogHandler) DeleteOldLogs(c *gin.Context) {
	daysStr := c.Param("days")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of days"})
		return
	}

	rowsAffected, err := h.repo.DeleteOlderThan(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Old logs deleted successfully", "rowsAffected": rowsAffected})
}

// LogActivity is a helper function to log an activity
func (h *ActivityLogHandler) LogActivity(c *gin.Context, action, entityType string, entityID interface{}, description string) {
	// Get the client IP address
	ipAddress := c.ClientIP()

	// For now, we'll use a default user ID of 1
	// In a real application, you would get this from the authenticated user
	userID := 1

	// Convert entityID to int if it's a string
	var entityIDInt int
	switch v := entityID.(type) {
	case int:
		entityIDInt = v
	case string:
		if v != "" {
			id, err := strconv.Atoi(v)
			if err == nil {
				entityIDInt = id
			}
		}
	default:
		entityIDInt = 0
	}

	// Log the activity asynchronously to not block the request
	go func() {
		h.repo.LogActivity(action, entityType, entityIDInt, description, userID, ipAddress)
	}()
}

// GetLogsCount handles GET /activity-logs/count
func (h *ActivityLogHandler) GetLogsCount(c *gin.Context) {
	filter := models.ActivityLogFilter{}

	count, err := h.repo.Count(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count activity logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetLogsByEntityTypeCount handles GET /activity-logs/entity-type/:entityType/count
func (h *ActivityLogHandler) GetLogsByEntityTypeCount(c *gin.Context) {
	entityType := c.Param("entityType")

	filter := models.ActivityLogFilter{
		EntityType: entityType,
	}

	count, err := h.repo.Count(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count activity logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetLogsByActionCount handles GET /activity-logs/action/:action/count
func (h *ActivityLogHandler) GetLogsByActionCount(c *gin.Context) {
	action := c.Param("action")

	filter := models.ActivityLogFilter{
		Action: action,
	}

	count, err := h.repo.Count(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count activity logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
