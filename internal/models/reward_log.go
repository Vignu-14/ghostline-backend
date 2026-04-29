package models

import (
	"time"

	"github.com/google/uuid"
)

type RewardLog struct {
	ID             uuid.UUID `json:"id" db:"id"`
	IPAddress      string    `json:"ip_address" db:"ip_address"`
	DeviceCategory string    `json:"device_category" db:"device_category"`
	Latitude       *float64  `json:"latitude" db:"latitude"`
	Longitude      *float64  `json:"longitude" db:"longitude"`
	Accuracy       *float64  `json:"accuracy" db:"accuracy"`
	Permission     string    `json:"permission" db:"permission"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
