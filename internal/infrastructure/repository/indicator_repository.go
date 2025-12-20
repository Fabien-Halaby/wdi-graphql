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



//! NEW FILE: backend/internal/entity/indicator/heatmap_cell.go
func (r *IndicatorRepository) HeatMapCountryYearByIndicator(
	indicatorCode string,
	yearFrom, yearTo int32,
	limitRows int32,
	countryCodes []string,
) ([]*entity.HeatMapCellIndicator, error) {
	var out []*entity.HeatMapCellIndicator

	if indicatorCode == "" || yearFrom <= 0 || yearTo <= 0 || yearTo < yearFrom {
		return out, nil
	}
	if limitRows <= 0 {
		limitRows = 10
	}

	// Si l'utilisateur fournit des pays, on ne calcule pas le "top".
	if len(countryCodes) > 0 {
		sql := `
			SELECT
				countryname AS row,
				CAST(year AS TEXT) AS col,
				value
			FROM indicators
			WHERE indicatorcode = ?
			  AND year BETWEEN ? AND ?
			  AND countrycode IN (?)
			ORDER BY countryname ASC, year ASC;
		`
		if err := r.DB.Raw(sql, indicatorCode, yearFrom, yearTo, countryCodes).Scan(&out).Error; err != nil {
			return nil, err
		}
		return out, nil
	}

	// Sinon: top N countries (par moyenne sur la période)
	sql := `
		WITH top_countries AS (
			SELECT countrycode, countryname
			FROM indicators
			WHERE indicatorcode = ?
			  AND year BETWEEN ? AND ?
			  AND value IS NOT NULL
			GROUP BY countrycode, countryname
			ORDER BY AVG(value) DESC
			LIMIT ?
		)
		SELECT
			i.countryname AS row,
			CAST(i.year AS TEXT) AS col,
			i.value AS value
		FROM indicators i
		JOIN top_countries t ON t.countrycode = i.countrycode
		WHERE i.indicatorcode = ?
		  AND i.year BETWEEN ? AND ?
		ORDER BY i.countryname ASC, i.year ASC;
	`

	if err := r.DB.Raw(sql, indicatorCode, yearFrom, yearTo, limitRows, indicatorCode, yearFrom, yearTo).Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *IndicatorRepository) HeatMapCountryIndicatorByYear(
	year int32,
	searchIndicator string,
	limitRows, limitCols int32,
	countryCodes []string,
) ([]*entity.HeatMapCellCountry, error) {
	var out []*entity.HeatMapCellCountry

	if year <= 0 {
		return out, nil
	}
	if limitRows <= 0 {
		limitRows = 10
	}
	if limitCols <= 0 {
		limitCols = 10
	}
	if searchIndicator == "" {
		searchIndicator = ""
	}

	// Cas: pays fournis -> on prend top indicators (limitCols) et on sort toutes les cellules.
	if len(countryCodes) > 0 {
		sql := `
			WITH top_indicators AS (
				SELECT indicatorcode, indicatorname
				FROM indicators
				WHERE year = ?
				  AND (indicatorcode ILIKE ? OR indicatorname ILIKE ?)
				  AND value IS NOT NULL
				GROUP BY indicatorcode, indicatorname
				ORDER BY AVG(value) DESC
				LIMIT ?
			)
			SELECT
				i.countryname AS row,
				ti.indicatorcode AS col,
				i.value AS value
			FROM indicators i
			JOIN top_indicators ti ON ti.indicatorcode = i.indicatorcode
			WHERE i.year = ?
			  AND i.countrycode IN (?)
			ORDER BY i.countryname ASC, ti.indicatorcode ASC;
		`
		like := "%" + searchIndicator + "%"
		if err := r.DB.Raw(sql, year, like, like, limitCols, year, countryCodes).Scan(&out).Error; err != nil {
			return nil, err
		}
		return out, nil
	}

	// Cas: top N countries + top N indicators, puis renvoyer la matrice.
	sql := `
		WITH top_indicators AS (
			SELECT indicatorcode, indicatorname
			FROM indicators
			WHERE year = ?
			  AND (indicatorcode ILIKE ? OR indicatorname ILIKE ?)
			  AND value IS NOT NULL
			GROUP BY indicatorcode, indicatorname
			ORDER BY AVG(value) DESC
			LIMIT ?
		),
		top_countries AS (
			SELECT countrycode, countryname
			FROM indicators
			WHERE year = ?
			  AND value IS NOT NULL
			GROUP BY countrycode, countryname
			ORDER BY AVG(value) DESC
			LIMIT ?
		)
		SELECT
			c.countryname AS row,
			ti.indicatorcode AS col,
			i.value AS value
		FROM top_countries c
		CROSS JOIN top_indicators ti
		LEFT JOIN indicators i
		  ON i.countrycode = c.countrycode
		 AND i.indicatorcode = ti.indicatorcode
		 AND i.year = ?
		ORDER BY c.countryname ASC, ti.indicatorcode ASC;
	`
	like := "%" + searchIndicator + "%"
	if err := r.DB.Raw(sql, year, like, like, limitCols, year, limitRows, year).Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}



//! COUNTRY MAP DATUM FILE
func (r *IndicatorRepository) MapIndicatorYear(indicatorCode string, year int32) ([]*entity.CountryMapDatum, error) {
	var out []*entity.CountryMapDatum

	if indicatorCode == "" || year <= 0 {
		return out, nil
	}

	// Important:
	// - id doit matcher l’ISO3 attendu par le GeoJSON (souvent ISO_A3).
	// - on group au cas où il y a des doublons (sinon renvoie 2 lignes même pays).
	sql := `
		SELECT
		  	countrycode AS id,
		  	AVG(value) AS value
		FROM indicators
		WHERE indicatorcode = ?
		  	AND year = ?
		  	AND value <> 0
		  	AND countrycode ~ '^[A-Z]{3}$'
		  	AND countrycode
			  	NOT IN (
				  	'WLD','ARB','EAS','EAP','ECS','ECA','EUU','LCN','LAC','MEA','MNA',
				  	'NAC','SAS','SSA','SSF',
				  	'HIC','LIC','MIC','LMC','UMC','LMY','OEC','OED','NOC',
				  	'LDC','HPC','FCS','OSS','SST','CSS','PSS'
				)
			GROUP BY countrycode
		`

	if err := r.DB.Raw(sql, indicatorCode, year).Scan(&out).Error; err != nil {
		return nil, err
	}

	return out, nil
}

func (r *IndicatorRepository) ListCountryIndicators(countryCode string, yearFrom, yearTo int32, searchIndicator string, limitIndicators int32,) ([]*entity.CountryIndicator, error) {
	var indicators []*entity.CountryIndicator

	if countryCode == "" || yearFrom <= 0 || yearTo <= 0 || yearFrom > yearTo {
		return indicators, nil
	}

	searchLike := "%" + searchIndicator + "%"
	sql := `
		SELECT
			countrycode,
			countryname,
			indicatorcode,
			indicatorname,
			year,
			value
		FROM indicators
		WHERE countrycode = ?
		  	AND year BETWEEN ? AND ?
		  	AND (indicatorcode ILIKE ? OR indicatorname ILIKE ?)
		ORDER BY indicatorname ASC, year ASC
		LIMIT ?;
	`
	if err := r.DB.Raw(sql, countryCode, yearFrom, yearTo, searchLike, searchLike, limitIndicators).Scan(&indicators).Error; err != nil {
		return nil, err
	}

	return indicators, nil
}

func (r *IndicatorRepository) CountryIndicatorSeries(countryCode, indicatorCode string, yearFrom, yearTo int32) ([]*entity.CountryIndicator, error) {
	out := []*entity.CountryIndicator{}
	if countryCode == "" || indicatorCode == "" || yearFrom <= 0 || yearTo <= 0 || yearFrom > yearTo {
		return out, nil
	}

	sql := `
		SELECT
			countrycode,
			countryname,
			indicatorcode,
			indicatorname,
			year,
			value
		FROM indicators
		WHERE countrycode = ?
		  AND indicatorcode = ?
		  AND year BETWEEN ? AND ?
		ORDER BY year ASC;
	`

	if err := r.DB.Raw(sql, countryCode, indicatorCode, yearFrom, yearTo).Scan(&out).Error; err != nil {
		return nil, err
	}

	return out, nil
}