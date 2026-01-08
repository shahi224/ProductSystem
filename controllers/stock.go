// controllers/stock.go
package controllers

import (
	"PRODUCT_SYSTEM/database"
	"PRODUCT_SYSTEM/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type StockRequest struct {
	SubVariantID string  `json:"sub_variant_id"`
	Qty          float64 `json:"qty"`
}

// stock in
func StockIn(c *gin.Context) {
	var req StockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var subVariant models.SubVariant
	if err := database.DB.First(&subVariant, "id = ?", req.SubVariantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SubVariant not found"})
		return
	}

	qty := decimal.NewFromFloat(req.Qty)
	subVariant.Stock = subVariant.Stock.Add(qty)

	if err := database.DB.Save(&subVariant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, "id = ?", subVariant.ProductID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Product not found"})
		return
	}

	var totalStock decimal.Decimal
	database.DB.Model(&models.SubVariant{}).
		Where("product_id = ?", subVariant.ProductID).
		Select("SUM(stock)").Scan(&totalStock)

	product.TotalStock = totalStock
	database.DB.Save(&product)

	transaction := models.StockTransaction{
		ID:              uuid.New().String(),
		ProductID:       subVariant.ProductID,
		SubVariantID:    subVariant.ID,
		Quantity:        qty,
		TransactionType: "IN",
		TransactionDate: time.Now(),
	}

	database.DB.Create(&transaction)

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock added successfully",
		"stock":   subVariant.Stock,
	})
}

// stock out
func StockOut(c *gin.Context) {
	var req StockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var subVariant models.SubVariant
	if err := database.DB.First(&subVariant, "id = ?", req.SubVariantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SubVariant not found"})
		return
	}

	qty := decimal.NewFromFloat(req.Qty)

	if subVariant.Stock.LessThan(qty) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	subVariant.Stock = subVariant.Stock.Sub(qty)

	if err := database.DB.Save(&subVariant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, "id = ?", subVariant.ProductID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Product not found"})
		return
	}

	var totalStock decimal.Decimal
	database.DB.Model(&models.SubVariant{}).
		Where("product_id = ?", subVariant.ProductID).
		Select("SUM(stock)").Scan(&totalStock)

	product.TotalStock = totalStock
	database.DB.Save(&product)

	transaction := models.StockTransaction{
		ID:              uuid.New().String(),
		ProductID:       subVariant.ProductID,
		SubVariantID:    subVariant.ID,
		Quantity:        qty,
		TransactionType: "OUT",
		TransactionDate: time.Now(),
	}

	database.DB.Create(&transaction)

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock removed successfully",
		"stock":   subVariant.Stock,
	})
}

// stock report
func StockReport(c *gin.Context) {
	type TransactionWithProduct struct {
		ID              string    `json:"id"`
		ProductID       string    `json:"product_id"`
		ProductName     string    `json:"product_name"`
		ProductCode     string    `json:"product_code"`
		SubVariantID    string    `json:"sub_variant_id"`
		Quantity        float64   `json:"quantity"`
		TransactionType string    `json:"transaction_type"`
		TransactionDate time.Time `json:"transaction_date"`
	}

	var transactions []TransactionWithProduct

	database.DB.Table("stock_transactions").
		Select(`
			stock_transaction.id,
			stock_transactions.product_id,
			products.product_name,
			products.product_code,
			stock_transactions.sub_variant_id,
			CAST(stock_transactions.quantity AS DECIMAL(10,2)) as quantity,
			stock_transactions.transaction_type,
			stock_transactions.transaction_date
		`).
		Joins("LEFT JOIN products ON products.id = stock_transactions.product_id").
		Order("stock_transactions.transaction_date DESC").
		Scan(&transactions)

	c.JSON(http.StatusOK, transactions)
}
