// models/stock.go
package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type StockTransaction struct {
	ID              string          `gorm:"type:char(36);primaryKey" json:"id"`
	ProductID       string          `gorm:"type:char(36)" json:"product_id"`
	SubVariantID    string          `gorm:"type:char(36)" json:"sub_variant_id"`
	Quantity        decimal.Decimal `gorm:"type:decimal(20,8)" json:"quantity"`
	TransactionType string          `gorm:"type:varchar(10)" json:"transaction_type"`
	TransactionDate time.Time       `json:"transaction_date"`
}
