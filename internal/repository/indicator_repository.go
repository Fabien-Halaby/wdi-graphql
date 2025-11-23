package repository

import (
	"wdi/internal/entity"

	"gorm.io/gorm"
)

type IndicatorRepository struct {
	DB *gorm.DB
}

func NewIndicatorRepository(db *gorm.DB) *IndicatorRepository {
	return &IndicatorRepository{DB: db}
}

func (r *IndicatorRepository) ListIndicators() ([]*entity.Indicator, error) {
	var indicators []*entity.Indicator
	if err := r.DB.Raw("SELECT * FROM Indicators").Find(&indicators).Error; err != nil {
		return nil, err
	}

	return indicators, nil
}
