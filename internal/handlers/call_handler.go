package handlers

import (
	"strings"

	"ghostline-backend/internal/config"
	"ghostline-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type CallHandler struct {
	config config.WebRTCConfig
}

func NewCallHandler(cfg config.WebRTCConfig) *CallHandler {
	return &CallHandler{config: cfg}
}

func (h *CallHandler) Config(c *fiber.Ctx) error {
	iceServers := make([]fiber.Map, 0, 2)

	if len(h.config.StunURLs) > 0 {
		iceServers = append(iceServers, fiber.Map{
			"urls": h.config.StunURLs,
		})
	}

	if len(h.config.TurnURLs) > 0 {
		iceServers = append(iceServers, fiber.Map{
			"urls":       h.config.TurnURLs,
			"username":   h.config.TurnUsername,
			"credential": h.config.TurnCredential,
		})
	}

	return utils.Success(c, fiber.StatusOK, "call configuration fetched successfully", fiber.Map{
		"ice_servers":      iceServers,
		"has_turn":         len(h.config.TurnURLs) > 0,
		"transport_policy": normalizeTransportPolicy(h.config.TransportPolicy),
	})
}

func normalizeTransportPolicy(value string) string {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		return config.DefaultWebRTCTransportPolicy
	}

	return trimmed
}
