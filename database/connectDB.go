package database

import (
	"PRODUCT_SYSTEM/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	_ = godotenv.Load()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",

		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db

	err = DB.AutoMigrate(
		&models.Product{},
		&models.Variant{},
		&models.VariantOption{},
		&models.SubVariant{},
		&models.StockTransaction{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("database connected & migrated successfully")
}
