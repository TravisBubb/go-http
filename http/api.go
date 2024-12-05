package http

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/TravisBubb/go-http/internal/tcp"
)

type endpoint struct {
    method HttpMethod
    path string
}

// Api represents an instance of a REST API that contain a set of endpoints with pre-defined handlers
type Api struct {
    endpoints map[endpoint]func(HttpRequest) HttpResponse
}

// Initialize a new Api
func CreateApi() *Api {
    var api Api
    api.endpoints = make(map[endpoint]func(HttpRequest) HttpResponse)
    return &api
}

// Map a new endpoint path template to a handler. This will return an error if the endpoint path already exists.
func (api *Api) Map(method HttpMethod, pathTemplate string, handler func(HttpRequest) HttpResponse) error {
    endpointKey := endpoint{method, pathTemplate}
    _, ok := api.endpoints[endpointKey]
    if ok {
        return fmt.Errorf("Path %s has already been registered", pathTemplate)
    }

    api.endpoints[endpointKey] = handler
    return nil
}

// Executes the API as a blocking call
func (api *Api) Run(port uint16) error {
    server := tcp.CreateServer("localhost", port, api.handleConnection)
    err := server.Start()
    if err != nil {
        log.Fatal("An unexpected error occurred attempting to start the API server:", err)
        return err
    }

    return nil
}

func (api *Api) handleConnection(ctx context.Context, connection net.Conn) {
	request, err := GetRequestFromConnection(ctx, connection)
	if err != nil {
		log.Println("Error occurred attempting to parse request:", err)
		return
	}

	log.Println("Request:", request)

    endpointKey := endpoint{request.Method, request.Path}

    // TODO: Handle path templates where route parameters are provided
    handleRequest, ok := api.endpoints[endpointKey]
    if !ok || handleRequest == nil {
        // TODO: Send back a properly formed HTTP response
        log.Printf("HTTP 404 Not Found - %s %s", request.Method.ToString(), request.Path)
        return
    }

    response := handleRequest(request)

    log.Println("Response:", response)
}
