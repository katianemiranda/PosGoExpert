package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

type Product struct {
	ID         int `gorm:"primaryKey"`
	Name       string
	Price      float64
	CategoryID int
	Category   Category
	gorm.Model
}

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/goexpert"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{}, &Category{})

	products := []Product{
		{Name: "Notebook", Price: 1000.0},
		{Name: "Mouse", Price: 50.0},
		{Name: "Keyboard", Price: 150.0},
	}
	db.Create((&products))

	// select one
	// var product Product2
	// db.First(&product, 1)
	// fmt.Println(product)
	// db.First(&product, "name", "Mouse")
	// fmt.Println(product)

	// select all
	// var products []Product2
	// db.Find(&products)
	// for _, products := range products {
	// 	fmt.Println(products)
	// }

	// select all
	// var products []Product2
	// db.Limit(2).Offset(2).Find(&products)
	// for _, product := range products {
	// 	fmt.Println(product)
	// }

	// where
	// var products []Product2
	// db.Where("price > ?", 100).Find(&products)
	// for _, product := range products {
	// 	fmt.Println(product)
	// }

	// var p Product2
	// db.First(&p, 1)
	// p.Name = "New Name"
	// db.Save(&p)

	// var p2 Product2
	// db.First(&p, 1)
	// fmt.Println(p2.Name)
	// db.Delete(&p)

}
