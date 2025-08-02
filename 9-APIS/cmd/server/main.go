package main

import (
	"katianemiranda/PosGoExpert/9-APIS/configs"
	"katianemiranda/PosGoExpert/9-APIS/internal/entity"
	"katianemiranda/PosGoExpert/9-APIS/internal/infra/database"
	"katianemiranda/PosGoExpert/9-APIS/internal/infra/webserver/handlers"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	http.HandleFunc("/products", productHandler.CreateProduct)
	http.ListenAndServe(":8000", nil)
}
