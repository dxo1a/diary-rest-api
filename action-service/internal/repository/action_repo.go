package repository

import (
	"action-service/internal/domain"
	"time"

	"gorm.io/gorm"
)

// управляет таблицей day_actions
type ActionRepo struct {
	db *gorm.DB
}

func NewActionRepo(db *gorm.DB) *ActionRepo {
	return &ActionRepo{db: db}
}

func (r *ActionRepo) ListByDate(date time.Time) ([]domain.DayAction, error) {
	var acts []domain.DayAction
	if err := r.db.Preload("Category").Where("date = ?", date).Find(&acts).Error; err != nil {
		return nil, err
	}
	return acts, nil
}

func (r *ActionRepo) Create(act *domain.DayAction) error {
	return r.db.Create(act).Error
}

func (r *ActionRepo) UpdateHours(id uint, hours float64) (domain.DayAction, error) {
	var act domain.DayAction
	if err := r.db.First(&act, id).Error; err != nil {
		return domain.DayAction{}, err
	}
	act.Hours = hours
	if err := r.db.Save(&act).Error; err != nil {
		return domain.DayAction{}, err
	}
	return act, nil
}

func (r *ActionRepo) Delete(id uint) error {
	return r.db.Delete(&domain.DayAction{}, id).Error
}
