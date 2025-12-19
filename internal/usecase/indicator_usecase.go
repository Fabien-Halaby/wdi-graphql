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
