// models/subvariant.go
package models

import "github.com/shopspring/decimal"

type SubVariant struct {
	ID        string          `gorm:"type:char(36);primaryKey" json:"id"`
	ProductID string          `gorm:"type:char(36)" json:"product_id"`
	OptionIDs string          `gorm:"type:text" json:"option_ids"`
	SKU       string          `json:"sku"`
	Stock     decimal.Decimal `gorm:"type:decimal(20,8);default:0" json:"stock"`
}
