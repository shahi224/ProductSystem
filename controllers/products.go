// controllers/product.go
package controllers

import (
	"PRODUCT_SYSTEM/database"
	"PRODUCT_SYSTEM/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateProduct - create a new product
func CreateProduct(c *gin.Context) {
	var req struct {
		ProductName  string `json:"product_name"`
		ProductImage string `json:"product_image"`
		CreatedUser  string `json:"created_user"`
		HSNCode      string `json:"hsn_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate unique product code
	productCode := "PROD-" + uuid.New().String()[:8]

	// Create Product
	product := models.Product{
		ID:           uuid.New().String(),
		ProductCode:  productCode,
		ProductName:  req.ProductName,
		ProductImage: req.ProductImage,
		CreatedDate:  time.Now(),
		CreatedUser:  req.CreatedUser,
		Active:       true,
		HSNCode:      req.HSNCode,
		TotalStock:   decimal.NewFromInt(0),
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// AUTO CREATE DEFAULT SUBVARIANT
	subVariant := models.SubVariant{
		ID:        uuid.New().String(),
		ProductID: product.ID,
		SKU:       product.ProductCode + "-DEFAULT",
		Stock:     decimal.NewFromInt(0),
	}

	if err := database.DB.Create(&subVariant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Product created successfully",
		"product":       product,
		"subvariant_id": subVariant.ID,
	})
}

// ListProducts with Pagination
func ListProducts(c *gin.Context) {
	// Get query parameters with defaults
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var products []models.Product
	var totalCount int64

	// Build query with search
	query := database.DB.Model(&models.Product{})

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("product_name LIKE ? OR product_code LIKE ? OR hsn_code LIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	// Get total count
	query.Count(&totalCount)

	// Get paginated products
	if err := query.Offset(offset).Limit(limit).
		Order("created_date DESC").
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total pages
	totalPages := (int(totalCount) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"products":    products,
		"page":        page,
		"limit":       limit,
		"total_items": totalCount,
		"total_pages": totalPages,
		"has_next":    page < totalPages,
		"has_prev":    page > 1,
	})
}

// GetProductWithSubVariants - get product with subvariants
func GetProductWithSubVariants(c *gin.Context) {
	productID := c.Param("id")

	var product models.Product
	if err := database.DB.First(&product, "id = ?", productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var subVariants []models.SubVariant
	if err := database.DB.Where("product_id = ?", productID).Find(&subVariants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product":     product,
		"subvariants": subVariants,
	})
}

// GetAllSubVariants - get all subvariants for dropdown
func GetAllSubVariants(c *gin.Context) {
	type SubVariantWithProduct struct {
		ID          string  `json:"id"`
		ProductID   string  `json:"product_id"`
		ProductName string  `json:"product_name"`
		ProductCode string  `json:"product_code"`
		SKU         string  `json:"sku"`
		Stock       float64 `json:"stock"`
	}

	var subVariants []SubVariantWithProduct

	// Join products and subvariants
	database.DB.Table("sub_variants").
		Select(`
            sub_variants.id,
            sub_variants.product_id,
            products.product_name,
            products.product_code,
            sub_variants.sku,
            CAST(sub_variants.stock AS DECIMAL(10,2)) as stock
        `).
		Joins("LEFT JOIN products ON products.id = sub_variants.product_id").
		Order("products.product_name, sub_variants.sku").
		Scan(&subVariants)

	c.JSON(http.StatusOK, gin.H{
		"subvariants": subVariants,
	})
}
