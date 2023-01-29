package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
    "os/signal"

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
    
    go func () {
        err := server.ListenAndServe()
        if err != nil {
            l.Fatal(err)
        }
    }()

    signalChannel := make(chan os.Signal)
    signal.Notify(signalChannel, os.Interrupt)
    signal.Notify(signalChannel, os.Kill)

    sig := <-signalChannel
    l.Println("Received terminate, graceful shutdown", sig)

    tc, _ := context.WithTimeout(context.Background(), 30 * time.Second)
    server.Shutdown(tc)
}
