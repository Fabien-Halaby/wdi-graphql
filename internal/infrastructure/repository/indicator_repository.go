package repository

import (
	entity "wdi/internal/entity/indicator"

	"gorm.io/gorm"
)

type IndicatorRepository struct {
	DB *gorm.DB
}

func NewIndicatorRepository(db *gorm.DB) *IndicatorRepository {
	return &IndicatorRepository{DB: db}
}

func (r *IndicatorRepository) ListIndicators(search string, limit int32, offset int32) ([]*entity.IndicatorList, error) {
	var indicators []*entity.IndicatorList
	if err := r.DB.Where("indicatorcode ILIKE ? OR indicatorname ILIKE ?", "%"+search+"%", "%"+search+"%").Limit(int(limit)).Offset(int(offset)).Find(&indicators).Error; err != nil {
		return nil, err
	}

	return indicators, nil
}
