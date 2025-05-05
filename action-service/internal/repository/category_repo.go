package repository

import (
	"action-service/internal/domain"

	"gorm.io/gorm"
)

// управляет таблицей action_categories
type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) List() ([]domain.ActionCategory, error) {
	var cats []domain.ActionCategory
	if err := r.db.Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

func (r *CategoryRepo) Create(name string) (domain.ActionCategory, error) {
	cat := domain.ActionCategory{Name: name}
	if err := r.db.Create(&cat).Error; err != nil {
		return domain.ActionCategory{}, err
	}
	return cat, nil
}

func (r *CategoryRepo) Update(id uint, name string) (domain.ActionCategory, error) {
	var cat domain.ActionCategory
	if err := r.db.First(&cat, id).Error; err != nil {
		return domain.ActionCategory{}, err
	}
	cat.Name = name
	if err := r.db.Save(&cat).Error; err != nil {
		return domain.ActionCategory{}, err
	}
	return cat, nil
}

func (r *CategoryRepo) Delete(id uint) error {
	return r.db.Delete(&domain.ActionCategory{}, id).Error
}
