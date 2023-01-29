package main

import ( 
    "net/http"
    "working/handlers"
    "log"
    "os"
)

func main() {
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
    hello := handlers.NewHello()

    serveMux := http.NewServeMux()
    serveMux.Handle("/", hh)

    http.ListenAndServe(":9090", nil)
}
