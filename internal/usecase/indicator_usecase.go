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

func (u *IndicatorUsecase) ListIndicators(search string, limit int32, offset int32) ([]*model.IndicatorList, error) {
	if limit == 0 || offset == 0 {
		limit = 50
		offset = 0
	}

	indicators, err := u.repo.ListIndicators(search, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*model.IndicatorList, 0, len(indicators))
	for _, i := range indicators {
		result = append(result, mapIndicatorEntityToModel(i))
	}

	return result, nil
}

func mapIndicatorEntityToModel(i *entity.IndicatorList) *model.IndicatorList {
	if i == nil {
		return nil
	}

	return &model.IndicatorList{
		Indicatorcode: i.IndicatorCode,
		Indicatorname: i.IndicatorName,
	}
}
