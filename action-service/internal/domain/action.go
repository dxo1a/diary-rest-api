package domain

import "time"

type DayAction struct {
	ID         uint           `gorm:"primaryKey"`
	Date       time.Time      `gorm:"not null"` // YYYY-MM-DD
	CategoryID uint           `gorm:"not null"` // ссылка на категорию
	Hours      float64        `gorm:"not null;check:hours>=0"`
	Category   ActionCategory `gorm:"foreignKey:CategoryID"`
}
