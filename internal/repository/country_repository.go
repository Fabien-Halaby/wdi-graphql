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

func (r *CountryRepository) FindAll(search *string, limit, offset int32) ([]*entity.Country, error) {
	var countries []*entity.Country
	db := r.DB
	if search != nil && *search != "" {
		pattern := "%" + *search + "%"
		db = db.Where(
			"COALESCE(ShortName, '') LIKE ? OR COALESCE(LongName, '') LIKE ? OR COALESCE(TableName, '') LIKE ?",
			pattern, pattern, pattern,
		)
	}
	if err := db.Limit(int(limit)).Offset(int(offset)).Order("ShortName ASC").Find(&countries).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("countries not found")
		}
		return nil, err
	}
	return countries, nil
}
