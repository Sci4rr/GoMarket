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
    err := json.NewDecoder(r.Body).Decode(&newToproduct)
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
    r.HandleFunc("/product/{id}", updateProductHandler).Methods("PUT")
    r.HandleFunc("/product/{id}", deleteProductHandler).Methods("DELETE")

    newProduct := Product{ID: "1", Name: "Test Product", Price: 9.99}
    productBytes, err := json.Marshal(newProduct)
    if err != nil {
        t.Fatalf("Failed to marshal new product: %v", err)
    }
    req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(productBytes))
    response := httptest.NewRecorder()
    r.ServeHTTP(response, req)

    if status := response.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
    }

    updatedProduct := Product{ID: "1", Name: "Updated Product", Price: 10.99}
    updatedProductBytes, err := json.Marshal(updatedProduct)
    if err != nil {
        t.Fatalf("Failed to marshal updated product: %v", err)
    }
    updateReq, _ := http.NewRequest("PUT", "/product/1", bytes.NewBuffer(updatedProductBytes))
    updateResponse := httptest.NewRecorder()
    r.ServeHTTP(updateResponse, updateReq)

    if status := updateResponse.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    deleteReq, _ := http.NewRequest("DELETE", "/product/1", nil)
    deleteResponse := httptest.NewRecorder()
    r.ServeHTTP(deleteResponse, deleteReq)

    if status := deleteResponse.Code; status != http.StatusNoContent {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
    }
}

func main() {
    fmt.Println("Run `go test` to test product management functionalities")
}