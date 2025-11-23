package usecase

import (
	"fmt"
	"wdi/graph/model"
	"wdi/internal/entity"
	"wdi/internal/repository"
)

type CountryUsecase struct {
	repo *repository.CountryRepository
}

func NewCountryUsecase(r *repository.CountryRepository) *CountryUsecase {
	return &CountryUsecase{repo: r}
}

func (u *CountryUsecase) GetAutocompleteCountries(prefix string, limit int) ([]*model.Country, error) {
	countries, err := u.repo.GetAutocompleteCountries(prefix, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Country, 0, len(countries))
	for _, c := range countries {
		result = append(result, mapCountryEntityToModel(c))
	}

	return result, nil
}

func (u *CountryUsecase) GetCountriesByCodes(codes []string) ([]*model.Country, error) {
	countries, err := u.repo.GetCountriesByCodes(codes)
	if err != nil {
		return nil, err
	}
	result := make([]*model.Country, 0, len(countries))
	for _, c := range countries {
		result = append(result, mapCountryEntityToModel(c))
	}

	return result, nil
}

func (u *CountryUsecase) CountByRegion() ([]*model.RegionCount, error) {
	rgcs, err := u.repo.CountByRegion()
	if err != nil {
		return nil, err
	}
	result := make([]*model.RegionCount, 0, len(rgcs))
	for _, r := range rgcs {
		if r.Region == "" {
			r.Region = "Uknown"
		}
		result = append(result, &model.RegionCount{
			Region: r.Region,
			Count:  int32(r.Count),
		})
	}

	return result, nil
}

func (u *CountryUsecase) CountCountries() (int, error) {
	return u.repo.CountCountries()
}

func (u *CountryUsecase) ListIncomeGroups() ([]string, error) {
	return u.repo.ListIncomeGroups()
}

func (u *CountryUsecase) ListRegions() ([]string, error) {
	regions, err := u.repo.ListRegions()
	if err != nil {
		return nil, err
	}

	for i := range regions {
		if regions[i] == "" {
			regions[i] = "Unkown"
		}
	}

	return regions, nil
}

func (u *CountryUsecase) GetAllCountries(search, region, incomeGroup *string, limit, offset int32) ([]*model.Country, error) {
	countries, err := u.repo.FindAll(search, region, incomeGroup, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Country, 0, len(countries))
	for _, c := range countries {
		result = append(result, mapCountryEntityToModel(c))
	}

	return result, nil
}

func (u *CountryUsecase) FindByCode(code string) (*model.Country, error) {
	if code == "" {
		return nil, fmt.Errorf("country code is required")
	}

	country, err := u.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	result := mapCountryEntityToModel(country)

	return result, nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func mapCountryEntityToModel(c *entity.Country) *model.Country {
	if c == nil {
		return nil
	}

	return &model.Country{
		CountryCode:  c.CountryCode,
		ShortName:    c.ShortName,
		TableName:    derefStr(c.TableNameDB),
		LongName:     derefStr(c.LongName),
		Alpha2Code:   derefStr(c.Alpha2Code),
		CurrencyUnit: derefStr(c.CurrencyUnit),
		SpecialNotes: derefStr(c.SpecialNotes),
		Region:       derefStr(c.Region),
		IncomeGroup:  derefStr(c.IncomeGroup),
	}
}
