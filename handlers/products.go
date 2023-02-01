// Package classification of Product API
// swagger:meta
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

// A list of products returns in the response
// swagger:response productsResponse
type productsResponseWrapper struct {
    // All products in the system
    // in: body
    Body []data.Product
}

// swagger:response noContent
type productsNoContent struct {
}

// swagger:parameters deleteProduct
type productIDParameterWrapper struct {
    // The id of the product to delete from the database
    // in: path
    // required: true
    ID int `json:"id"`
}

type Products struct{
    l* log.Logger
}

func NewProducts(l *log.Logger) *Products{
    return &Products{l}
}


// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//   200: productsResponse

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request){
    listProd := data.GetProducts()
    err := listProd.ToJSON(rw)
    if err != nil {
        http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
    }
}

// AddProduct adds a product from the database
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
    p.l.Println("Handle POST Product")

    prod := r.Context().Value(KeyProduct{}).(data.Product)
    data.AddProduct(&prod)
}


// UpdateProducts updates a product from the database
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

// swagger:route DELETE /products/{id} products deleteProduct 
// Deletes a product from database
// responses:
//   201: noContent 

// DeleteProduct deletes a product from the database
func (p * Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, _ := strconv.Atoi(vars["id"])

    p.l.Println("Handle DELETE Product", id)

    //err := data.DeleteProduct(id)

    /*if err == data.ErrProductNotFound {
        http.Error(rw, "Product not found", http.StatusNotFound)
        return
    }

    if err != nil {
        http.Error(rw, "Product not found", http.StatusInternalServerError)
    }*/
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
