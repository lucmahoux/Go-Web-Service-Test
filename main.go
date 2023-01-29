package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"github.com/lucmahoux/go_http_test/handlers"
)

func main() {
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
    hello := handlers.NewHello(l)
    goodbye := handlers.NewGoodbye(l)

    serveMux := http.NewServeMux()
    serveMux.Handle("/", hello)
    serveMux.Handle("/goodbye", goodbye)

    server := &http.Server{
        Addr: ":9090",
        Handler: serveMux,
        IdleTimeout: 120*time.Second,
        ReadTimeout: 1*time.Second,
        WriteTimeout: 1*time.Second,
    }

    server.ListenAndServe()
}
