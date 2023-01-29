package handlers

import (
	"log"
	"net/http"
	"github.com/lucmahoux/go_http_test/data"
)

type Products struct{
    l* log.Logger
}

func NewProducts(l *log.Logger) *Products{
    return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, h *http.Request) {
    listProd := data.GetProducts()
    err := listProd.ToJSON(rw)
    if err != nil {
        http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
    }
}
