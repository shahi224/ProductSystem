package models

type Variant struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	ProductID string `gorm:"type:char(36)"`
	Name      string

	Options []VariantOption `gorm:"foreignKey:VariantID"`
}

type VariantOption struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	VariantID string `gorm:"type:char(36)"`
	Value     string
}
