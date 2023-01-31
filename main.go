package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
    "os/signal"
    "github.com/gorilla/mux"
	"github.com/lucmahoux/go_http_test/handlers"
)

func main() {
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
    
    // create the handlers
    //hello := handlers.NewHello(l)
    //goodbye := handlers.NewGoodbye(l)
    productHandler := handlers.NewProducts(l)

    // create a new serve mux and register the handlers
    serveMux := mux.NewRouter()

    getRouter := serveMux.Methods(http.MethodGet).Subrouter()
    getRouter.HandleFunc("/", productHandler.GetProducts)

    putRouter := serveMux.Methods(http.MethodPut).Subrouter()
    putRouter.HandleFunc("/{id:[0-9]+}", productHandler.UpdateProducts)

    postRouter := serveMux.Methods(http.MethodPost).Subrouter()
    postRouter.HandleFunc("/", productHandler.AddProduct)

    // create a new server
    server := &http.Server{
        Addr:               ":9090",            // configure the bind address
        Handler:            serveMux,           // set the default handler
        IdleTimeout:        120*time.Second,
        ReadTimeout:        1*time.Second,
        WriteTimeout:       1*time.Second,
    }
    
    // start the server
    go func () {
        err := server.ListenAndServe()
        if err != nil {
            l.Fatal(err)
        }
    }()

    // trap sigterm or interrupt and gracefully shutdown the server
    signalChannel := make(chan os.Signal)
    signal.Notify(signalChannel, os.Interrupt)
    signal.Notify(signalChannel, os.Kill)

    sig := <-signalChannel
    l.Println("Received terminate, graceful shutdown", sig)

    tc, _ := context.WithTimeout(context.Background(), 30 * time.Second)
    server.Shutdown(tc)
}
