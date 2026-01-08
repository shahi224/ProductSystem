package main

import (
	"PRODUCT_SYSTEM/database"
	"PRODUCT_SYSTEM/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	fmt.Println("database connected & migrated successfully")

	r := gin.Default()

	//  HTML templates
	r.LoadHTMLGlob("templates/*.html")

	//  Static files
	r.Static("/static", "./static")

	routes.RegisterRoutes(r)

	// Serve HTML pages
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/create", func(c *gin.Context) {
		c.HTML(200, "create.html", nil)
	})

	r.GET("/list", func(c *gin.Context) {
		c.HTML(200, "list.html", nil)
	})

	r.GET("/stock", func(c *gin.Context) {
		c.HTML(200, "stock.html", nil)
	})

	r.GET("/report", func(c *gin.Context) {
		c.HTML(200, "report.html", nil)
	})

	r.Run(":8080")
}
