package entity

type Indicator struct {
	CountryCode   string  `gorm:"column:CountryCode"`
	CountryName   string  `gorm:"column:CountryName"`
	IndicatorCode string  `gorm:"column:IndicatorCode"`
	IndicatorName string  `gorm:"column:IndicatorName"`
	Year          int     `gorm:"column:Year"`
	Value         float64 `gorm:"column:Value"`
}

func (Indicator) TableName() string {
	return "Indicators"
}
