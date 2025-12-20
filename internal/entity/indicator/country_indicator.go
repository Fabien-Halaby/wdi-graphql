package entity

type CountryIndicator struct {
	CountryCode   string `gorm:"column:countrycode" json:"countrycode"`
	CountryName   string	`gorm:"column:countryname" json:"countryname"`
	IndicatorCode string `gorm:"column:indicatorcode" json:"indicatorcode"`
	IndicatorName string `gorm:"column:indicatorname" json:"indicatorname"`
	Year          int32  `gorm:"column:year" json:"year"`
	Value         float64 `gorm:"column:value" json:"value"`
}

func (CountryIndicator) TableName() string {
	return "indicators"
}