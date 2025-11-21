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

func (u *CountryUsecase) GetAllCountries() ([]*model.Country, error) {
	countries, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]*model.Country, 0, len(countries))
	for _, c := range countries {
		result = append(result, &model.Country{
			CountryCode:  c.CountryCode,
			ShortName:    c.ShortName,
			TableName:    *c.DTableName,
			LongName:     *c.LongName,
			Alpha2Code:   *c.Alpha2Code,
			CurrencyUnit: *c.CurrencyUnit,
			SpecialNotes: *c.SpecialNotes,
			Region:       *c.Region,
			IncomeGroup:  *c.IncomeGroup,
		})
	}

	return result, nil
}

// func (u *CountryUsecase) GetCountryByCode(code string) (*entity.Country, error) {
// 	return u.repo.FindByCode(code)
// }

// func (u *CountryUsecase) GetRegionCounts() ([]repository.RegionCount, error) {
// 	return u.repo.RegionCounts()
// }
