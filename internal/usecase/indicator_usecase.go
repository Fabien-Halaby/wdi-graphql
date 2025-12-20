package usecase

import (
	entity "wdi/internal/entity/indicator"
	"wdi/internal/infrastructure/repository"
	"wdi/internal/interface/graph/model"
)

type IndicatorUsecase struct {
	repo *repository.IndicatorRepository
}

func NewIndicatorUsecase(repo *repository.IndicatorRepository) *IndicatorUsecase {
	return &IndicatorUsecase{repo: repo}
}

func (u *IndicatorUsecase) ListIndicators(search string, limit int32, offset int32) ([]*model.Indicator, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	indicators, err := u.repo.ListIndicators(search, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Indicator, 0, len(indicators))
	for _, i := range indicators {
		result = append(result, mapIndicatorEntityToModel(i))
	}

	return result, nil
}

func (u *IndicatorUsecase) ListIndicatorYears() ([]*model.IndicatorYear, error) {
	years, err := u.repo.ListIndicatorYears()
	if err != nil {
		return nil, err
	}

	result := make([]*model.IndicatorYear, 0, len(years))
	for _, y := range years {
		result = append(result, &model.IndicatorYear{
			Year: y.Year,
		})
	}

	return result, nil
}

func (u *IndicatorUsecase) ListIndicatorBarRace(indicatorCode string, limitCountries int32) ([]*model.IndicatorBarRaceRow, error) {
	rows, err := u.repo.ListIndicatorBarRace(indicatorCode, limitCountries)
	if err != nil {
		return nil, err
	}

	result := make([]*model.IndicatorBarRaceRow, 0, len(rows))
	for _, r := range rows {
		result = append(result, &model.IndicatorBarRaceRow{
			Year:    r.Year,
			Country: r.Country,
			Value:   r.Value,
		})
	}
	return result, nil
}


func mapIndicatorEntityToModel(i *entity.Indicator) *model.Indicator {
	if i == nil {
		return nil
	}

	return &model.Indicator{
		Code: i.Code,
		Name: i.Name,
	}
}


//! NEW FILE
func (u *IndicatorUsecase) HeatMapCountryYearByIndicator(indicatorCode string, yearFrom, yearTo int32, limitRows int32, countryCodes []string) ([]*model.HeatMapCell, error) {
	cells, err := u.repo.HeatMapCountryYearByIndicator(indicatorCode, yearFrom, yearTo, limitRows, countryCodes)
	if err != nil {
		return nil, err
	}

	out := make([]*model.HeatMapCell, 0, len(cells))
	for _, c := range cells {
		out = append(out, &model.HeatMapCell{
			Row:   c.Row,
			Col:   c.Col,
			Value: c.Value, // gqlgen gère Float nullable via *float64
		})
	}
	return out, nil
}

func (u *IndicatorUsecase) HeatMapCountryIndicatorByYear(year int32, searchIndicator string, limitRows, limitCols int32, countryCodes []string) ([]*model.HeatMapCell, error) {
	cells, err := u.repo.HeatMapCountryIndicatorByYear(year, searchIndicator, limitRows, limitCols, countryCodes)
	if err != nil {
		return nil, err
	}

	out := make([]*model.HeatMapCell, 0, len(cells))
	for _, c := range cells {
		out = append(out, &model.HeatMapCell{
			Row:   c.Row,
			Col:   c.Col,
			Value: c.Value,
		})
	}
	return out, nil
}


//! COUNTRY MAP DATUM
func (u *IndicatorUsecase) MapIndicatorYear(indicatorCode string, year int32) ([]*model.CountryMapDatum, error) {
	rows, err := u.repo.MapIndicatorYear(indicatorCode, year)
	if err != nil {
		return nil, err
	}

	out := make([]*model.CountryMapDatum, 0, len(rows))
	for _, r := range rows {
		out = append(out, &model.CountryMapDatum{
			ID:    r.ID,
			Value: r.Value,
		})
	}
	return out, nil
}


func (u *IndicatorUsecase) ListCountryIndicators(countryCode string, yearFrom, yearTo int32, searchIndicator string, limitIndicators int32) ([]*model.CountryIndicatorValue, error) {
	indicators, err := u.repo.ListCountryIndicators(countryCode, yearFrom, yearTo, searchIndicator, limitIndicators)
	if err != nil {
		return nil, err
	}

	result := make([]*model.CountryIndicatorValue, 0, len(indicators))
	for _, i := range indicators {
		result = append(result, &model.CountryIndicatorValue{
			Countrycode:   i.CountryCode,
			Countryname:   i.CountryName,
			Indicatorcode: i.IndicatorCode,
			Indicatorname: i.IndicatorName,
			Year:          i.Year,
			Value:         i.Value,
		})
	}

	return result, nil
}

func (u *IndicatorUsecase) CountryIndicatorSeries(countryCode string, indicatorCode string, yearFrom int32, yearTo int32) ([]*model.CountryIndicatorValue, error) {
	indicators, err := u.repo.CountryIndicatorSeries(countryCode, indicatorCode, yearFrom, yearTo)
	if err != nil {
		return nil, err
	}

	result := make([]*model.CountryIndicatorValue, 0, len(indicators))
	for _, i := range indicators {
		result = append(result, &model.CountryIndicatorValue{
			Countrycode:   i.CountryCode,
			Countryname:   i.CountryName,
			Indicatorcode: i.IndicatorCode,
			Indicatorname: i.IndicatorName,
			Year:          i.Year,
			Value:         i.Value,
		})
	}

	return result, nil
}