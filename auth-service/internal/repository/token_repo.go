package repository

import (
	"auth-service/internal/domain"
	"time"

	"gorm.io/gorm"
)

type TokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) Create(rt *domain.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *TokenRepo) Find(token string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	if err := r.db.Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *TokenRepo) Revoke(id uint) error {
	return r.db.Model(&domain.RefreshToken{}).Where("id = ?", id).Update("revoked", true).Error
}

func (r *TokenRepo) DeleteExpired(now time.Time) error {
	return r.db.Where("expires_at < ?", now).Delete(&domain.RefreshToken{}).Error
}
