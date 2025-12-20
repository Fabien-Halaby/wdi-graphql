package entity

type CountryMapDatum struct {
	ID    string  `gorm:"column:id" json:"id"`
	Value float64 `gorm:"column:value" json:"value"`
}

func (CountryMapDatum) TableName() string {
	return "indicators"
}