package usecase

import (
	"wdi/graph/model"
	"wdi/internal/entity"
	"wdi/internal/repository"
)

type IndicatorUsecase struct {
	repo *repository.IndicatorRepository
}

func NewIndicatorUsecase(repo *repository.IndicatorRepository) *IndicatorUsecase {
	return &IndicatorUsecase{repo: repo}
}

func (u *IndicatorUsecase) ListIndicators() ([]*model.Indicator, error) {
	indicators, err := u.repo.ListIndicators()
	if err != nil {
		return nil, err
	}

	result := make([]*model.Indicator, 0, len(indicators))
	for _, i := range indicators {
		result = append(result, mapIndicatorEntityToModel(i))
	}

	return result, nil
}

func mapIndicatorEntityToModel(i *entity.Indicator) *model.Indicator {
	if i == nil {
		return nil
	}

	return &model.Indicator{
		CountryCode: i.CountryCode,
		CountryName: i.CountryName,
		Code:        i.IndicatorCode,
		Name:        i.IndicatorName,
		Year:        int32(i.Year),
		Value:       i.Value,
	}
}
