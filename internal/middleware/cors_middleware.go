package middleware

import (
	"log/slog"
	"strings"

	"ghostline-backend/internal/config"

	"github.com/gofiber/fiber/v2"
	fibercors "github.com/gofiber/fiber/v2/middleware/cors"
)

// NewCORS creates a CORS middleware that supports multiple comma-separated origins
func NewCORS(cfg config.CORSConfig) fiber.Handler {
	// Parse comma-separated origins
	allowedOrigins := parseOrigins(cfg.AllowedOrigin)
	slog.Info("CORS configuration loaded", "allowed_origins", allowedOrigins)

	return fibercors.New(fibercors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return IsOriginAllowed(origin, allowedOrigins)
		},
		AllowCredentials: cfg.AllowCredentials,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		ExposeHeaders:    "X-Request-ID, Retry-After",
		MaxAge:           86400, // 24 hours for preflight cache
	})
}

// IsOriginAllowed checks if a given origin is in the allowed list
func IsOriginAllowed(origin string, allowedOrigins []string) bool {
	// Allow requests with no origin (same-origin requests, mobile apps)
	if origin == "" {
		return true
	}

	// Check if origin is in the allowed list
	for _, allowed := range allowedOrigins {
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}

	// Log rejected origins for debugging
	slog.Warn("Origin rejected", "origin", origin, "allowed", allowedOrigins)
	return false
}

// parseOrigins splits comma-separated origin strings into a slice
func parseOrigins(originStr string) []string {
	if strings.TrimSpace(originStr) == "" {
		return []string{config.DefaultAllowedOrigin} // Use default from constants
	}

	origins := strings.Split(originStr, ",")
	result := make([]string, 0, len(origins))
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
