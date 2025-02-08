package handlers

import (
	"ecommerce/database/dao"
	"ecommerce/models"
	"ecommerce/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// CreateProduct handles creating a new product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	product.ID = utils.NewID()
	product.CreatedDate = time.Now()
	product.UpdatedDate = time.Now()

	if err := dao.CreateProduct(&product); err != nil {
		log.Printf("unable to create product : %s", err)
		http.Error(w, "Unable to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GetProducts handles fetching all products
func GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := dao.GetProducts()
	if err != nil {
		http.Error(w, "Unable to fetch products", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}
