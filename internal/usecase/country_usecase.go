package usecase

import (
	"fmt"
	"wdi/graph/model"
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
		result = append(result, &model.Country{
			CountryCode:  c.CountryCode,
			ShortName:    c.ShortName,
			TableName:    derefStr(c.TableNameDb),
			LongName:     derefStr(c.LongName),
			Alpha2Code:   derefStr(c.Alpha2Code),
			CurrencyUnit: derefStr(c.CurrencyUnit),
			SpecialNotes: derefStr(c.SpecialNotes),
			Region:       derefStr(c.Region),
			IncomeGroup:  derefStr(c.IncomeGroup),
		})
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
		result = append(result, &model.Country{
			CountryCode:  c.CountryCode,
			ShortName:    c.ShortName,
			TableName:    derefStr(c.TableNameDb),
			LongName:     derefStr(c.LongName),
			Alpha2Code:   derefStr(c.Alpha2Code),
			CurrencyUnit: derefStr(c.CurrencyUnit),
			SpecialNotes: derefStr(c.SpecialNotes),
			Region:       derefStr(c.Region),
			IncomeGroup:  derefStr(c.IncomeGroup),
		})
	}

	return result, nil
}

func (u *CountryUsecase) GetNumberOfCountryPerRegion() ([]*model.RegionCount, error) {
	rgcs, err := u.repo.GetNumberOfCountryPerRegion()
	if err != nil {
		return nil, err
	}
	result := make([]*model.RegionCount, 0, len(rgcs))
	for _, r := range rgcs {
		result = append(result, &model.RegionCount{
			Region: r.Region,
			Count:  int32(r.Count),
		})
	}

	return result, nil
}

func (u *CountryUsecase) GetNumberOfAllCountry() (int, error) {
	return u.repo.GetNumberOfAllCountry()
}

func (u *CountryUsecase) GetAllIncomeGroups() ([]string, error) {
	return u.repo.GetAllIncomeGroups()
}

func (u *CountryUsecase) GetAllRegions() ([]string, error) {
	return u.repo.GetAllRegions()
}

func (u *CountryUsecase) GetAllCountries(search, region, incomeGroup *string, limit, offset int32) ([]*model.Country, error) {
	countries, err := u.repo.FindAll(search, region, incomeGroup, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Country, 0, len(countries))
	for _, c := range countries {
		result = append(result, &model.Country{
			CountryCode:  c.CountryCode,
			ShortName:    c.ShortName,
			TableName:    derefStr(c.TableNameDb),
			LongName:     derefStr(c.LongName),
			Alpha2Code:   derefStr(c.Alpha2Code),
			CurrencyUnit: derefStr(c.CurrencyUnit),
			SpecialNotes: derefStr(c.SpecialNotes),
			Region:       derefStr(c.Region),
			IncomeGroup:  derefStr(c.IncomeGroup),
		})
	}

	return result, nil
}

func (u *CountryUsecase) FindByCode(code string) (*model.Country, error) {
	if code == "" {
		return nil, fmt.Errorf("failed to gagner l'ordinateur")
	}

	country, err := u.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	result := &model.Country{
		CountryCode:  country.CountryCode,
		ShortName:    country.ShortName,
		TableName:    derefStr(country.TableNameDb),
		LongName:     derefStr(country.LongName),
		Alpha2Code:   derefStr(country.Alpha2Code),
		CurrencyUnit: derefStr(country.CurrencyUnit),
		SpecialNotes: derefStr(country.SpecialNotes),
		Region:       derefStr(country.Region),
		IncomeGroup:  derefStr(country.IncomeGroup),
	}

	return result, nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
