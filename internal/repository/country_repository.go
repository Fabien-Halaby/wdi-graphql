package repository

import (
	"errors"
	"wdi/internal/entity"

	"gorm.io/gorm"
)

type CountryRepository struct {
	DB *gorm.DB
}

func NewCountryRepository(db *gorm.DB) *CountryRepository {
	return &CountryRepository{DB: db}
}

func (r *CountryRepository) FindAll(limit, offset int32) ([]*entity.Country, error) {
	var countries []*entity.Country
	if err := r.DB.Limit(int(limit)).Offset(int(offset)).Order("ShortName ASC").Find(&countries).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("countries not found")
		}
		return nil, err
	}

	return countries, nil
}
