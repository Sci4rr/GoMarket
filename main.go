package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
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

func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getUser(w http.ResponseWriter, r *http.Request) {
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

func createTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transactions = append(transactions, transaction)
	json.NewEncoder(w).Encode(transaction)
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != os.Getenv("API_KEY") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func main() {
	initEnv()

	router := mux.NewRouter()

	router.HandleFunc("/api/products", getProducts).Methods("GET")
	router.HandleFunc("/api/user/{username}", getUser).Methods("GET")
	router.HandleFunc("/api/transaction", createTransaction).Methods("POST")

	router.Use(authMiddleware)

	log.Fatal(http.ListenAndServe(":8080", router))
}