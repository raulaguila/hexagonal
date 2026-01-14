package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/raulaguila/go-api/internal/adapter/driver/rest/presenter"
	"github.com/raulaguila/go-api/internal/app"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	app *app.Application
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(router fiber.Router, application *app.Application) {
	handler := &HealthHandler{
		app: application,
	}
	router.Get("", handler.healthCheck).Name("Root")
	router.Get("/health", handler.detailedHealth).Name("Health")
}

// HealthStatus represents the health check response
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version,omitempty"`
	Uptime    string                 `json:"uptime,omitempty"`
	Checks    map[string]CheckResult `json:"checks,omitempty"`
}

// CheckResult represents the result of a health check
type CheckResult struct {
	Status   string `json:"status"`
	Duration string `json:"duration,omitempty"`
	Message  string `json:"message,omitempty"`
}

var startTime = time.Now()

// healthCheck godoc
// @Summary      Ping Pong
// @Description  Simple health check endpoint
// @Tags         Health
// @Produce      json
// @Success      200  {object}   map[string]string
// @Router       / [get]
func (h *HealthHandler) healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":    "ok",
		"timestamp": time.Now(),
	})
}

// detailedHealth godoc
// @Summary      Detailed Health Check
// @Description  Returns detailed health status including database and storage checks
// @Tags         Health
// @Produce      json
// @Success      200  {object}   HealthStatus
// @Failure      503  {object}   HealthStatus
// @Router       /health [get]
func (h *HealthHandler) detailedHealth(c *fiber.Ctx) error {
	checks := make(map[string]CheckResult)
	overallStatus := "healthy"

	// Check database
	dbCheck := h.checkDatabase()
	checks["database"] = dbCheck
	if dbCheck.Status != "up" {
		overallStatus = "unhealthy"
	}

	response := HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   h.app.Version(),
		Uptime:    time.Since(startTime).Round(time.Second).String(),
		Checks:    checks,
	}

	statusCode := fiber.StatusOK
	if overallStatus != "healthy" {
		statusCode = fiber.StatusServiceUnavailable
	}

	return presenter.New(c, statusCode, overallStatus, response)
}

// checkDatabase checks the database connection
func (h *HealthHandler) checkDatabase() CheckResult {
	start := time.Now()

	// Try to ping the database through repository
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Simple check: try to count users with limit 0
	_, err := h.app.Repositories.User.Count(ctx, nil)

	duration := time.Since(start)

	if err != nil {
		return CheckResult{
			Status:   "down",
			Duration: duration.String(),
			Message:  err.Error(),
		}
	}

	return CheckResult{
		Status:   "up",
		Duration: duration.String(),
	}
}
