package repository

import (
	"errors"
	"fmt"
	"wdi/internal/entity"

	"gorm.io/gorm"
)

type CountryRepository struct {
	DB *gorm.DB
}

func NewCountryRepository(db *gorm.DB) *CountryRepository {
	return &CountryRepository{DB: db}
}

func (r *CountryRepository) GetAutocompleteCountries(prefix string, limit int) ([]*entity.Country, error) {
	var countries []*entity.Country
	prefix = prefix + "%"
	if err := r.DB.Raw("SELECT * FROM Country WHERE ShortName LIKE ? OR LongName LIKE ? OR TableName LIKE ? LIMIT ?", prefix, prefix, prefix, limit).Limit(limit).Find(&countries).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no result matches")
		}
		return nil, err
	}

	return countries, nil
}

func (r *CountryRepository) GetCountriesByCodes(codes []string) ([]*entity.Country, error) {
	var countries []*entity.Country
	if err := r.DB.Raw("SELECT * FROM Country WHERE CountryCode IN ?", codes).Find(&countries).Error; err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *CountryRepository) GetNumberOfAllCountry() (int, error) {
	var nb int
	if err := r.DB.Raw("SELECT COUNT(*) FROM Country").First(&nb).Error; err != nil {
		return 0, err
	}

	return nb, nil
}

func (r *CountryRepository) GetNumberOfCountryPerRegion() ([]*entity.RegionCount, error) {
	var rgc []*entity.RegionCount
	if err := r.DB.Raw("SELECT Region, COUNT(*) AS Count FROM Country WHERE Region IS NOT NULL GROUP BY Region").Find(&rgc).Error; err != nil {
		return nil, err
	}

	return rgc, nil
}

func (r *CountryRepository) GetAllIncomeGroups() ([]string, error) {
	var incomeGroups []string
	if err := r.DB.Raw("SELECT DISTINCT IncomeGroup FROM Country WHERE IncomeGroup IS NOT NULL AND IncomeGroup <> '' ORDER BY IncomeGroup").First(&incomeGroups).Error; err != nil {
		return nil, err
	}

	return incomeGroups, nil
}

func (r *CountryRepository) GetAllRegions() ([]string, error) {
	var region []string
	if err := r.DB.Raw("SELECT DISTINCT Region FROM Country WHERE Region IS NOT NULL AND Region <> '' ORDER BY Region ASC").Find(&region).Error; err != nil {
		return nil, err
	}

	return region, nil
}

func (r *CountryRepository) FindByCode(code string) (*entity.Country, error) {
	var country *entity.Country
	if err := r.DB.Where("CountryCode = ?", code).Limit(1).First(&country).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("country not found")
		}
		return nil, err
	}

	return country, nil
}

func (r *CountryRepository) FindAll(search, region, incomeGroup *string, limit, offset int32) ([]*entity.Country, error) {
	var countries []*entity.Country
	db := r.DB

	if region != nil && *region != "" {
		db = db.Where("Region = ?", *region)
	}
	if incomeGroup != nil && *incomeGroup != "" {
		db = db.Where("IncomeGroup = ?", *incomeGroup)
	}
	if search != nil && *search != "" {
		pattern := "%" + *search + "%"
		db = db.Where(
			"COALESCE(ShortName, '') LIKE ? OR COALESCE(LongName, '') LIKE ? OR COALESCE(TableName, '') LIKE ?",
			pattern, pattern, pattern,
		)
	}

	if err := db.
		Limit(int(limit)).
		Offset(int(offset)).
		Order("ShortName ASC").
		Find(&countries).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("countries not found")
		}
		return nil, err
	}

	return countries, nil
}
