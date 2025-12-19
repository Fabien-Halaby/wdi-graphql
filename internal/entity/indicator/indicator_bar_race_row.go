package entity

type IndicatorBarRaceRow struct {
	Year    int32   `gorm:"column:year" json:"year"`
	Country string  `gorm:"column:country" json:"country"`
	Value   float64 `gorm:"column:value" json:"value"`
}
