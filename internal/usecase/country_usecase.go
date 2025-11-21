package usecase

import (
	"wdi/graph/model"
	"wdi/internal/repository"
)

type CountryUsecase struct {
	repo *repository.CountryRepository
}

func NewCountryUsecase(r *repository.CountryRepository) *CountryUsecase {
	return &CountryUsecase{repo: r}
}

func (u *CountryUsecase) GetAllCountries(limit, offset int32) ([]*model.Country, error) {
	countries, err := u.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Country, 0, len(countries))
	for _, c := range countries {
		result = append(result, &model.Country{
			CountryCode:  c.CountryCode,
			ShortName:    c.ShortName,
			TableName:    DerefStr(c.TableNameDb),
			LongName:     DerefStr(c.LongName),
			Alpha2Code:   DerefStr(c.Alpha2Code),
			CurrencyUnit: DerefStr(c.CurrencyUnit),
			SpecialNotes: DerefStr(c.SpecialNotes),
			Region:       DerefStr(c.Region),
			IncomeGroup:  DerefStr(c.IncomeGroup),
		})
	}

	return result, nil
}

func DerefStr(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
