package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Hello is a simple handler
type Hello struct {
    log *log.Logger
}

// NewHello creates a new hello handler with the given logger
func NewHello(l *log.Logger) *Hello {
    return &Hello{l}
}

// ServeHTTP implements the go http.Handler interface
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    h.log.Println("Hello World")
    
    // read the body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(rw, "Oops", http.StatusBadRequest)
        return
    }
    
    // write the response
    fmt.Fprintf(rw, "Hello %s", body)
}
