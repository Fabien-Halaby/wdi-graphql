package repository

import (
	"wdi/internal/entity"

	"gorm.io/gorm"
)

type CountryRepository struct {
	DB *gorm.DB
}

func NewCountryRepository(db *gorm.DB) *CountryRepository {
	return &CountryRepository{DB: db}
}

// Simple GORM pour lister tous les pays
func (r *CountryRepository) FindAll() ([]*entity.Country, error) {
	var countries []*entity.Country
	if err := r.DB.Find(&countries).Error; err != nil {
		return nil, err
	}

	return countries, nil
}

// Simple GORM pour trouver un pays par code
// func (r *CountryRepository) FindByCode(code string) (*entity.Country, error) {
// 	var c entity.Country
// 	if err := r.DB.Where("CountryCode = ?", code).First(&c).Error; err != nil {
// 		return nil, err
// 	}
// 	return &c, nil
// }

// // SQL pur : nombre de pays par région
// type RegionCount struct {
// 	Region string
// 	Count  int
// }

// func (r *CountryRepository) RegionCounts() ([]RegionCount, error) {
// 	var rows []RegionCount
// 	sql := `
//         SELECT Region, COUNT(*) AS Count
//         FROM Country
//         WHERE Region IS NOT NULL
//         GROUP BY Region
//     `
// 	if err := r.DB.Raw(sql).Scan(&rows).Error; err != nil {
// 		return nil, err
// 	}
// 	return rows, nil
// }
