package handlers

import (
	"ghostline-backend/internal/models"
	"ghostline-backend/internal/repositories"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type RewardHandler struct {
	repo *repositories.RewardRepository
}

func NewRewardHandler(repo *repositories.RewardRepository) *RewardHandler {
	return &RewardHandler{repo: repo}
}

func (h *RewardHandler) LogLocation(c *fiber.Ctx) error {
	slog.Info("Reward location logging request received")

	var body struct {
		DeviceCategory string   `json:"device_category"`
		Permission     string   `json:"permission"`
		Latitude       *float64 `json:"latitude"`
		Longitude      *float64 `json:"longitude"`
		Accuracy       *float64 `json:"accuracy"`
	}

	if err := c.BodyParser(&body); err != nil {
		slog.Error("Failed to parse reward location body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	ip := c.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.IP()
	}

	rewardLog := &models.RewardLog{
		IPAddress:      ip,
		DeviceCategory: body.DeviceCategory,
		Permission:     body.Permission,
		Latitude:       body.Latitude,
		Longitude:      body.Longitude,
		Accuracy:       body.Accuracy,
	}

	slog.Info("Attempting to insert reward log", 
		"ip", ip, 
		"permission", body.Permission,
		"lat", body.Latitude,
		"lng", body.Longitude,
	)

	if err := h.repo.Create(c.Context(), rewardLog); err != nil {
		slog.Error("Failed to insert reward location to database", "error", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "error",
			"error":   err.Error(),
		})
	}

	slog.Info("Reward location logged successfully", "id", rewardLog.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id": rewardLog.ID,
		},
	})
}
