package entity

type IndicatorYear struct {
	Year int32 `gorm:"column:year" json:"year"`
}

func (IndicatorYear) TableName() string {
	return "mv_year_list"
}
