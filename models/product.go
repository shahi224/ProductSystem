// models/product.go
package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID           string          `gorm:"type:char(36);primaryKey" json:"id"`
	ProductID    int64           `gorm:"unique;autoIncrement" json:"product_id"`
	ProductCode  string          `gorm:"unique" json:"product_code"`
	ProductName  string          `json:"product_name"`
	ProductImage string          `json:"product_image"`
	CreatedDate  time.Time       `gorm:"autoCreateTime" json:"created_date"`
	UpdatedDate  *time.Time      `gorm:"autoUpdateTime" json:"updated_date"`
	CreatedUser  string          `gorm:"type:char(36)" json:"created_user"`
	IsFavourite  bool            `json:"is_favourite"`
	Active       bool            `json:"active"`
	HSNCode      string          `json:"hsn_code"`
	TotalStock   decimal.Decimal `gorm:"type:decimal(20,8);default:0" json:"total_stock"`

	Variants []Variant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
}
