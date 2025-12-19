package entity

type Indicator struct {
	Code string `gorm:"column:indicatorcode" json:"code"`
	Name string `gorm:"column:indicatorname" json:"name"`
}

func (Indicator) TableName() string {
	return "mv_indicator_list"
}
