// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// Basepath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta

package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lucmahoux/go_http_test/data"
)

type Products struct{
    l* log.Logger
}

func NewProducts(l *log.Logger) *Products{
    return &Products{l}
}


func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request){
    listProd := data.GetProducts()
    err := listProd.ToJSON(rw)
    if err != nil {
        http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
    }
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
    p.l.Println("Handle POST Product")

    prod := r.Context().Value(KeyProduct{}).(data.Product)
    data.AddProduct(&prod)
}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(rw, "Unable to convert id", http.StatusBadRequest)
        return
    }

    p.l.Println("Handle PUT Product")
    prod := r.Context().Value(KeyProduct{}).(data.Product)

    err = data.UpdateProduct(id, &prod)
    if err == data.ErrProductNotFound {
        http.Error(rw, "Product not found", http.StatusNotFound)
        return
    }

    if err != nil{
        http.Error(rw, "Product not found", http.StatusInternalServerError)
        return
    }
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
    return http.HandlerFunc(func (rw http.ResponseWriter, r *http.Request) {
        prod := data.Product{}

        err := prod.FromJSON(r.Body)
        if err != nil {
            p.l.Println("[ERROR] deserializing product", err)
            http.Error(rw, "Error reading product", http.StatusBadRequest)
            return
        }
        
        //validate the product
        err = prod.Validate()
        if err != nil {
            p.l.Println("[ERROR] validate product", err)
            http.Error(
                rw, 
                fmt.Sprintf("Error validating product: %s", err),
                http.StatusBadRequest,
            )
            return
        }

        // add the product to the context
        ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
        req := r.WithContext(ctx)

        // call the next handler, which can be another middleware in the chain,
        // or the final handler
        next.ServeHTTP(rw, req)
    })
}
