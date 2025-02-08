package handlers

import (
	"ecommerce/database"
	"ecommerce/database/dao"
	"ecommerce/models"
	"ecommerce/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// AddToCart handles adding a product to a user's cart
func AddToCart(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart
	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	cart.ID = utils.NewID()
	cart.CreatedDate = time.Now()

	if err := dao.AddProductToCart(&cart); err != nil {
		log.Printf("unable to add product to cart : %s", err)
		http.Error(w, "Unable to add product to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

// GetCart handles fetching the cart for a user
func GetCartItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	var cartItems models.CartItems
	cartItems.UserID = userID

	items, err := dao.GetCartItems(database.DB, userID)
	if err != nil {
		log.Printf("unable to fetch cart items of userID : %s, err : %s", userID, err)
		http.Error(w, "Unable to fetch cart", http.StatusInternalServerError)
		return
	}

	cartItems.Products = items

	json.NewEncoder(w).Encode(cartItems)
}
