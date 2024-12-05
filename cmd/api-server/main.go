package main

import (
    "log"
    "github.com/TravisBubb/go-http/http"
)

func main() {
    log.Println("Creating api...")
    api := http.CreateApi()

    log.Println("Registering endpoints...")
    _ = api.Map(http.GET, "/v1/products", getProducts)

    log.Println("Starting api...")
    err := api.Run(8080)
    if err != nil {
        log.Fatal("An error occurred attempting to run api server", err)
    }
}

func getProducts(request http.HttpRequest) http.HttpResponse {
    return http.HttpResponse{
       StatusCode: http.OK, 
       Protocol: request.Protocol,
       Content: "{\"Products\": []}",
    }
}
