package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Products []Product
}

type Product struct {
	ID         int `gorm:"primaryKey"`
	Name       string
	Price      float64
	CategoryID int
	Category   Category
	gorm.Model
}

// type SerialNumber struct {
// 	ID        int `gorm:"primaryKey"`
// 	Number    string
// 	ProductID int
// }

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//db.AutoMigrate(&Product{}, &Category{}, &SerialNumber{})
	db.AutoMigrate(&Product{}, &Category{})

	//create category
	// category := Category{Name: "Cozinha"}
	// db.Create(&category)

	// create product
	// db.Create(&Product{
	// 	Name:       "Notebook",
	// 	Price:      1000.0,
	// 	CategoryID: category.ID,
	// })

	// create product
	// db.Create(&Product{
	// 	Name:       "Mouse",
	// 	Price:      1000.0,
	// 	CategoryID: category.ID,
	// })

	// create product
	// db.Create(&Product{
	// 	Name:       "Panela",
	// 	Price:      99.0,
	// 	CategoryID: category.ID,
	// })

	// //create serial number
	// db.Create(&SerialNumber{
	// 	Number:    "123456",
	// 	ProductID: 1,
	// })

	var products []Product
	db.Preload("Category").Find(&products)
	//db.Find(&products)
	for _, product := range products {
		fmt.Println(product.Name, product.Category.Name)
	}

	var cateories []Category
	err = db.Model(&Category{}).Preload("Products").Find(&cateories).Error
	if err != nil {
		panic(err)
	}
	for _, category := range cateories {
		fmt.Println(category.Name)
		for _, product := range category.Products {
			fmt.Println(" - ", product.Name, category.Name)
		}
	}

}
