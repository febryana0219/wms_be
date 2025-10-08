package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
	"wms-be/domain/models"
	"wms-be/infrastructure/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Store(userID uuid.UUID, token string, expiresAt time.Time) error
	FindByToken(token string) (*models.RefreshToken, error)
	Revoke(token string) error
	DeleteExpired() error
}

type refreshTokenRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepo{db: database.GetDB()}
}

// hashToken: simpan hash token ke DB, bukan token asli
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Store simpan refresh token ke DB
func (r *refreshTokenRepo) Store(userID uuid.UUID, token string, expiresAt time.Time) error {
	rt := &models.RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(token),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return r.db.Create(rt).Error
}

// FindByToken cek token valid dan belum revoked
func (r *refreshTokenRepo) FindByToken(token string) (*models.RefreshToken, error) {
	tokenHash := hashToken(token)
	var rt models.RefreshToken

	err := r.db.Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&rt).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("refresh token not found or expired")
	}
	return &rt, err
}

// Revoke tandai token sebagai revoked
func (r *refreshTokenRepo) Revoke(token string) error {
	tokenHash := hashToken(token)
	now := time.Now()

	return r.db.Model(&models.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Updates(map[string]interface{}{
			"revoked_at": now,
			"updated_at": now,
		}).Error
}

// DeleteExpired bersihkan token yang sudah lewat masa berlaku
func (r *refreshTokenRepo) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}
