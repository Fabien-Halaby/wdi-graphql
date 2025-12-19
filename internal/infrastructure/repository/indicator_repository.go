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

func (r *IndicatorRepository) ListIndicators(search string, limit, offset int32) ([]*entity.Indicator, error) {
	var indicators []*entity.Indicator

	q := r.DB.Model(&entity.Indicator{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("indicatorcode ILIKE ? OR indicatorname ILIKE ?", like, like)
	}

	if err := q.Limit(int(limit)).Offset(int(offset)).Find(&indicators).Error; err != nil {
		return nil, err
	}

	return indicators, nil
}

func (r *IndicatorRepository) ListIndicatorYears() ([]*entity.IndicatorYear, error) {
	var years []*entity.IndicatorYear

	if err := r.DB.Model(&entity.IndicatorYear{}).Find(&years).Error; err != nil {
		return nil, err
	}

	return years, nil
}

func (r *IndicatorRepository) ListIndicatorBarRace(indicatorCode string, limitCountries int32) ([]*entity.IndicatorBarRaceRow, error) {
	var out []*entity.IndicatorBarRaceRow

	if indicatorCode == "" {
		return out, nil
	}
	if limitCountries <= 0 {
		limitCountries = 10
	}

	sql := `
		SELECT year, country, value
		FROM (
			SELECT
				year,
				countryname AS country,
				value,
				ROW_NUMBER() OVER (PARTITION BY year ORDER BY value DESC) AS rn
			FROM indicators
			WHERE indicatorcode = ?
			  AND value IS NOT NULL
		) ranked
		WHERE rn <= ?
		ORDER BY year ASC, value DESC;
	`

	if err := r.DB.Raw(sql, indicatorCode, limitCountries).Scan(&out).Error; err != nil {
		return nil, err
	}

	return out, nil
}
