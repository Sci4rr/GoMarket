package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Transaction struct {
	ProductID string  `json:"productID"`
	Username  string  `json:"username"`
	Amount    float64 `json:"amount"`
}

var products []Product
var users []User
var transactions []Transaction

var (
	productCache     []byte
	productCacheLock sync.RWMutex
	cacheNeedsUpdate = true
)

func LoadEnvironmentVariables() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productCacheLock.RLock()
	if cacheNeedsUpdate {
		productCacheLock.RUnlock()
		productCacheLock.Lock()
		if cacheNeedsUpdate {
			var err error
			productCache, err = json.Marshal(products)
			if err != nil {
				http.Error(w, "Failed to encode products", http.StatusInternalServerError)
				return
			}
			cacheNeedsUpdate = false
		}
		productCacheLock.Unlock()
		productCacheLock.RLock()
	}
	w.Write(productCache)
	productCacheLock.RUnlock()
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		http.Error(w, "Username is missing in parameters", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.Username == username {
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.NotFound(w, r)
}

func HandleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transactions = append(transactions, transaction)
	cacheNeedsUpdate = true
	json.NewEncoder(w).Encode(transaction)
}

func WithAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != os.Getenv("API_KEY") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func main() {
	LoadEnvironmentVariables()

	router := mux.NewRouter()

	router.HandleFunc("/api/products", HandleGetProducts).Methods("GET")
	router.HandleFunc("/api/user/{username}", HandleGetUser).Methods("GET")
	router.HandleFunc("/api/transaction", HandleCreateTransaction).Methods("POST")

	router.Use(WithAuthMiddleware)

	log.Fatal(http.ListenAndServe(":8080", router))
}