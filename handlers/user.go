package handlers

import (
	"ecommerce/database/dao"
	"ecommerce/kafka"
	"ecommerce/models" // A utility package for encryption
	"ecommerce/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	producer *kafka.Producer
}

func NewUser(producer *kafka.Producer) *User {
	return &User{
		producer: producer,
	}
}

// CreateUser handles user registration
func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user.ID = utils.NewID()
	user.CreatedDate = time.Now()
	user.UpdatedDate = time.Now()

	// Hashing the password
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	err := dao.CreateUser(&user)
	if err != nil {
		http.Error(w, "Unable to create user in db", http.StatusBadRequest)
		return
	}

	userInfo := models.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		Email:     user.Email,
	}
	u.producer.PublishUserAccountCreated(&userInfo)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate user credentials
	user, err := dao.GetUserByEmail(loginData.Email)
	if err != nil || !utils.CheckPasswordHash(loginData.Password, user.Password) {
		log.Printf("db errror is %s", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT Token for authentication
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		log.Fatal(err)
	}

	if len(token) == 0 {
		log.Fatal("invalid token!")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// GetUser handles fetching user information
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("user id is : %s", id)

	user, err := dao.GetUserByID(id)
	if err != nil {
		log.Printf("db error is : %s", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
