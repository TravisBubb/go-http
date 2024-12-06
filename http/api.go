package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

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

    endpointKey := endpoint{request.Method, request.Path}

    responseHeaders := make(map[string]string)
    responseHeaders["Content-Type"] = "application/json"

    // TODO: Handle path templates where route parameters are provided
    handleRequest, ok := api.endpoints[endpointKey]
    if !ok || handleRequest == nil {
        response := CreateResponse(NotFound, "", responseHeaders, request.Protocol)
        responseString := formatHttpResponse(response)
        _, err = connection.Write([]byte(responseString))
        if err != nil {
            log.Println("Error occurred attempting to send response:", err)
        }
        return
    }

    response := handleRequest(request)
    response.Headers = responseHeaders
    responseString := formatHttpResponse(response)
    _, err = connection.Write([]byte(responseString))
    if err != nil {
        log.Println("Error occurred attempting to send response:", err)
    }
}

func formatHttpResponse(response HttpResponse) string {
    var s strings.Builder
    s.WriteString(fmt.Sprintf("%s %d %s\n", response.Protocol, response.StatusCode, response.StatusCode.ToString()))

    for k, v := range response.Headers {
        s.WriteString(fmt.Sprintf("%s: %s\n", k, v))
    }

    s.WriteString(fmt.Sprintf("\n%s", response.Content))

    return s.String() 
}
