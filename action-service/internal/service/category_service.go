package service

import (
	"action-service/internal/domain"
	"action-service/internal/repository"
)

// бизнес-логика для категорий
type CategoryService struct {
	repo *repository.CategoryRepo
}

func NewCategoryService(repo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List() ([]domain.ActionCategory, error) {
	return s.repo.List()
}

func (s *CategoryService) Create(name string) (domain.ActionCategory, error) {
	return s.repo.Create(name)
}

func (s *CategoryService) Update(id uint, name string) (domain.ActionCategory, error) {
	return s.repo.Update(id, name)
}

func (s *CategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}
