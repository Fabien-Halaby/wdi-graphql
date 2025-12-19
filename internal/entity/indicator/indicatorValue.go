package entity

type IndicatorValue struct {
	ID           int64     `gorm:"id" json:"id"`
	Code 	   string    `gorm:"indicatorcode" json:"code"`
	Name 	   string    `gorm:"indicatorname" json:"name"`
	Value 	   float64   `gorm:"value" json:"value"`
	Year 	   int       `gorm:"year" json:"year"`
	CountryCode  string    `gorm:"countrycode" json:"country_code"`
	CountryName  string    `gorm:"countryname" json:"country_name"`
}

func (IndicatorValue) TableName() string {
	return "indicators"
}