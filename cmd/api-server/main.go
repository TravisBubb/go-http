package main

import (
	"github.com/TravisBubb/go-http/http"
	"log"
)

func main() {
	log.Println("Creating api...")
	api := http.CreateApi()

	log.Println("Registering endpoints...")
	_ = api.Map(http.GET, "/v1/products", func(c *http.Context) {
		c.Ok("this was successful...")
	})

	_ = api.Map(http.POST, "/v1/products", func(c *http.Context) {
		var request postProductRequest
		err := c.BindRequest(&request)
		if err != nil {
			log.Println("Error binding request:", err)
			return
		}

		c.Ok("ID " + request.Id)
	})

	log.Println("Starting api...")
	err := api.Run(8080)
	if err != nil {
		log.Fatal("An error occurred attempting to run api server", err)
	}
}

type postProductRequest struct {
	Id string `json:"id"`
}
