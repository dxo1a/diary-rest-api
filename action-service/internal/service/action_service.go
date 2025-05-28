package service

import (
	"action-service/internal/domain"
	"action-service/internal/repository"
	"time"
)

// бизнес-логика для действий по дням
type ActionService struct {
	repo *repository.ActionRepo
}

func NewActionService(repo *repository.ActionRepo) *ActionService {
	return &ActionService{repo: repo}
}

func (s *ActionService) List(date time.Time) ([]domain.DayAction, error) {
	return s.repo.ListByDate(date)
}

func (s *ActionService) Add(date time.Time, categoryID uint, hours float64) (domain.DayAction, error) {
	act := domain.DayAction{
		Date:       date,
		CategoryID: categoryID,
		Hours:      hours,
	}
	err := s.repo.Create(&act)
	return act, err
}

func (s *ActionService) Update(id uint, hours float64) (domain.DayAction, error) {
	return s.repo.UpdateHours(id, hours)
}

func (s *ActionService) Delete(id uint) error {
	return s.repo.Delete(id)
}
