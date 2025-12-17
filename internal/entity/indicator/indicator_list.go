package entity

type IndicatorList struct {
	IndicatorCode string `gorm:"column:indicatorcode" json:"indicatorcode"`
	IndicatorName string `gorm:"column:indicatorname" json:"indicatorname"`
}

func (IndicatorList) TableName() string {
	return "mv_indicator_list"
}
