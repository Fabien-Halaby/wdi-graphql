package entity

type HeatMapCellCountry struct {
	Row   string   `gorm:"column:row" json:"row"`
	Col   string   `gorm:"column:col" json:"col"`
	Value *float64 `gorm:"column:value" json:"value"`
}

type HeatMapCellIndicator struct {
	Row   string   `gorm:"column:row" json:"row"`
	Col   string   `gorm:"column:col" json:"col"`
	Value *float64 `gorm:"column:value" json:"value"`
}

func (HeatMapCellCountry) TableName() string {
	return "indicators"
}

func (HeatMapCellIndicator) TableName() string {
	return "indicators"
}