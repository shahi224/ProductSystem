Product Management System 🏪

A comprehensive inventory and stock management system built with Go (Gin) backend and HTML/CSS/JavaScript frontend.

✨ Features

📦 Product Management

Create products with auto-generated subvariants
View all products with pagination (10 products per page)
Search products by name, code, or HSN
Product statistics dashboard
📊 Stock Management

Stock IN/OUT operations
Real-time stock tracking
Subvariant-based inventory control
Stock transaction history
📈 Reporting

Complete stock transaction report
Filter by transaction type (IN/OUT)
Date-based filtering
Export functionality
🛠️ Tech Stack

Backend

Go (Golang) - Programming language
Gin - Web framework
GORM - ORM for MySQL
MySQL - Database
uuid - Unique identifier generation
decimal - Precise decimal calculations
Frontend

HTML5 - Markup
CSS3 - Styling
Vanilla JavaScript - Frontend logic
Responsive Design - Works on all devices


📁 Project Structure

PRODUCT_SYSTEM/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Dependency checksums
├── .env.example           # Environment template (✅ Safe to commit)
├── .gitignore            # Git ignore rules
├── README.md             # This file
│
├── controllers/           # HTTP controllers
│   ├── product.go        # Product CRUD operations
│   └── stock.go          # Stock management
│
├── database/             # Database configuration
│   └── database.go      # MySQL connection & migration
│
├── models/               # Database models
│   ├── product.go       # Product model
│   ├── stock.go         # Stock transaction model
│   ├── subvariant.go    # Subvariant model
│   └── variant.go       # Variant model
│
├── routes/               # URL routing
│   └── routes.go        # Route definitions
│
├── templates/            # HTML templates
│   ├── index.html       # Home page
│   ├── create.html      # Create product
│   ├── list.html        # Product list
│   ├── stock.html       # Stock management
│   └── report.html      # Stock report
│
└── static/               # Static assets
    ├── app.js           # Frontend JavaScript
    └── style.css        # Stylesheet
