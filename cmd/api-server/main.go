package main

import (
    "log"

    "github.com/TravisBubb/go-http/http"
)

func main() {
    server := http.CreateServer("localhost", 8080)
    err := server.Start()
    if err != nil {
        log.Fatal("An unexpected error occurred attempting to start server:", err)
    }
}
