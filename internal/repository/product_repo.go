package repository

import "github.com/imran4u/Oolio/internal/model"

var products = []model.Product{
	{"1", "Chicken Waffle", 10.5, "Waffle"},
	{"2", "Veg Burger", 8.5, "Burger"},
	{"3", "French Fries", 5.0, "Snack"},
}

func GetAllProducts() []model.Product {
	return products
}

func GetProductByID(id string) (*model.Product, bool) {

	for _, p := range products {
		if p.ID == id {
			return &p, true
		}
	}

	return nil, false
}
