package main

import ( 
    "net/http"
    "github.com/lucmahoux/go_http_test/handlers"
    "log"
    "os"
)

func main() {
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
    hello := handlers.NewHello(l)
    goodbye := handlers.NewGoodbye(l)

    serveMux := http.NewServeMux()
    serveMux.Handle("/", hello)
    serveMux.Handle("/goodbye", goodbye)

    http.ListenAndServe(":9090", serveMux)
}
