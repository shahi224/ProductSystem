package routes

import (
	"PRODUCT_SYSTEM/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.POST("/products", controllers.CreateProduct)
	r.GET("/products", controllers.ListProducts)
	r.GET("/products/:id/subvariants", controllers.GetProductWithSubVariants)
	r.GET("/subvariants/all", controllers.GetAllSubVariants) // Add this line

	r.POST("/stock/in", controllers.StockIn)
	r.POST("/stock/out", controllers.StockOut)
	r.GET("/stock/report", controllers.StockReport)
}
