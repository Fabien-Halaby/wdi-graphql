package resolvers

import (
	"wdi/internal/usecase"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	IndicatorUC *usecase.IndicatorUsecase
	DB          *gorm.DB
}
