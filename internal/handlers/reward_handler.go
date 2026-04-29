package handlers

import (
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/repositories"
	"log"

	"github.com/gofiber/fiber/v2"
)

type RewardHandler struct {
	repo *repositories.RewardRepository
}

func NewRewardHandler(repo *repositories.RewardRepository) *RewardHandler {
	return &RewardHandler{repo: repo}
}

func (h *RewardHandler) LogLocation(c *fiber.Ctx) error {
	var body struct {
		DeviceCategory string   `json:"device_category"`
		Permission     string   `json:"permission"`
		Latitude       *float64 `json:"latitude"`
		Longitude      *float64 `json:"longitude"`
		Accuracy       *float64 `json:"accuracy"`
	}

	if err := c.BodyParser(&body); err != nil {
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

	if err := h.repo.Create(c.Context(), rewardLog); err != nil {
		log.Printf("Failed to log reward location: %v", err)
		// We return 200 even on error to ensure frontend flow continues
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
