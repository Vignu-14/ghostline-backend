package repositories

import (
	"anonymous-communication/backend/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RewardRepository struct {
	db *pgxpool.Pool
}

func NewRewardRepository(db *pgxpool.Pool) *RewardRepository {
	return &RewardRepository{db: db}
}

func (r *RewardRepository) Create(ctx context.Context, log *models.RewardLog) error {
	query := `
		INSERT INTO reward_logs (
			ip_address, device_category, latitude, longitude, accuracy, permission
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		log.IPAddress,
		log.DeviceCategory,
		log.Latitude,
		log.Longitude,
		log.Accuracy,
		log.Permission,
	).Scan(&log.ID, &log.CreatedAt)
}
