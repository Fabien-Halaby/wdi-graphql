package entity

type Country struct {
	CountryCode  string  `gorm:"column:CountryCode"`
	ShortName    string  `gorm:"column:ShortName"`
	LongName     *string `gorm:"column:LongName"`
	Region       *string `gorm:"column:Region"`
	DTableName   *string `gorm:"column:TableName"`
	Alpha2Code   *string `gorm:"column:Alpha2Code"`
	SpecialNotes *string `gorm:"column:SpecialNotes"`
	IncomeGroup  *string `gorm:"column:IncomeGroup"`
	CurrencyUnit *string `gorm:"column:CurrencyUnit"`
}

func (Country) TableName() string {
	return "Country"
}
