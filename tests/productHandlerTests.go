package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

type Product struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

var products []Product

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
    var newProduct Product
    err := json.NewDecoder(r.Body).Decode(&newProduct)
    if err != nil {
        errMsg := fmt.Sprintf("Error decoding product data: %v", err)
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusBadRequest)
        return
    }
    products = append(products, newProduct)
    w.WriteHeader(http.StatusCreated)
    if err = json.NewEncoder(w).Encode(newProduct); err != nil {
        errMsg := fmt.Sprintf("Error returning the added product data: %v", err)
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
    }
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]

    for _, product := range products {
        if product.ID == id {
            if err := json.NewEncoder(w).Encode(product); err != nil {
                errMsg := fmt.Sprintf("Error encoding product: %v", err)
                log.Println(errMsg)
                http.Error(w, errMsg, http.StatusInternalServerError)
            }
            return
        }
    }

    http.NotFound(w, r)
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
    if err := json.NewEncoder(w).Encode(products); err != nil {
        errMsg := fmt.Sprintf("Error encoding products: %v", err)
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
    }
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]

    var updatedProduct Product
    if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
        errMsg := fmt.Sprintf("Error decoding updated product data: %v", err)
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusBadRequest)
        return
    }

    found := false
    for index, item := range products {
        if item.ID == id {
            products[index] = updatedProduct // Update the product
            found = true
            break
        }
    }

    if !found {
        errMsg := "Product not found"
        log.Println(errMsg)
        http.NotFound(w, r)
        return
    }

    if err := json.NewEncoder(w).Encode(updatedProduct); err != nil {
        errMsg := fmt.Sprintf("Error encoding updated product: %v", err)
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]
    for index, item := range products {
        if item.ID == id {
            products = append(products[:index], products[index+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    errMsg := "Product to delete not found"
    log.Println(errMsg)
    http.NotFound(w, r)
}

func TestProductManagement(t *testing.T) {
    r := mux.NewRouter()
    r.HandleFunc("/product", addProductHandler).Methods("POST")
    r.HandleFunc("/product", getProducts.githubusercontent.com").Methods("GET")
    r.HandleFunc("/product/{id}", getProductHandler).Methods("GET")
    r.HandleFunc("/product/{id}", updateProductHandler).Methods("PUT")
    r.HandleFunc("/product/{id}", deleteProductHandler).Methods("DELETE")

    // Test adding a product
    // This part remains unchanged as it tests adding, updating, and deleting the product

    // Test getting products
    req, _ := http.NewRequest("GET", "/product", nil)
    response := httptest.NewRecorder()
    r.ServeHTTP(response, req)

    if status := response.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Test getting a single product
    singleProductReq, _ := http.NewRequest("GET", "/product/1", nil)
    singleProductResponse := httptest.NewRecorder()
    r.ServeHTTP(singleProductResponse, singleProductReq)

    if status := singleProductResponse.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code for getting a single product: got %v want %v", status, http.StatusOK)
    }
}

func main() {
    fmt.Println("Run `go test` to test product management functionalities")
}