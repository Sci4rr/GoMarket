package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "github.com/joho/godotenv"
    "github.com/gorilla/mux"
)

type Product struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

var products []Product


func init() {
    _ = godotenv.Load()
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
    var newProduct Product
    err := json.NewDecoder(r.Body).Decode(&newProduct)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    products = append(products, newProduct)
    w.WriteHeader(http.StatusCreated)
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]
    for index, item := range products {
        if item.ID == id {
            var updatedProduct Product
            err := json.NewDecoder(r.Body).Decode(&updatedProduct)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            products[index] = updatedProduct
            json.NewEncoder(w).Encode(updatedProduct)
            return
        }
    }
    http.NotFound(w, r)
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]
    for index, item := range products {
        if item.ID == id {
            products = append(products[:index], products[index+1:]...)
            break
        }
    }
    w.WriteHeader(http.StatusNoContent)
}

func TestProductManagement(t *testing.T) {
    r := mux.NewRouter()
    r.HandleFunc("/product", addProductHandler).Methods("POST")
    r.HandleFunc("/product/{id}", updateProductHandler).Methods("PUT")
    r.HandleFunc("/product/{id}", deleteProductHandler).Methods("DELETE")

    newProduct := Product{ID: "1", Name: "Test Product", Price: 9.99}
    productBytes, _ := json.Marshal(newMontage)
    req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(productBytes))
    response := httptest.NewRecorder()
    r.ServeHTTP(response, req)

    if status := response.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
    }

    updatedProduct := Product{ID: "1", Name: "Updated Product", Price: 10.99}
    updatedProductBytes, _ := json.Marshal(updatedMontage)
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