package http

import(
    "fmt"
)

// Enum for HTTP methods
type HttpMethod int

const (
	GET HttpMethod = iota
	POST
	PUT
	DELETE
)

func (method HttpMethod) ToString() string {
    switch method {
        case GET:
            return "GET"
        case POST:
            return "POST"
        case PUT:
            return "PUT"
        case DELETE:
            return "DELETE"
        default:
            panic(fmt.Sprintf("HTTP Method %d not supported", method))
    }
}

// Enum for HTTP status codes
type HttpStatusCode int

const (
	OK                  HttpStatusCode = 200
	Created             HttpStatusCode = 201
	BadRequest          HttpStatusCode = 400
	Unauthorized        HttpStatusCode = 401
	NotFound            HttpStatusCode = 404
	InternalServerError HttpStatusCode = 500
)

func (statusCode HttpStatusCode) ToString() string {
    switch statusCode {
        case OK:
            return "OK"
        case Created:
            return "Created"
        case BadRequest:
            return "Bad Request"
        case Unauthorized:
            return "Unauthorized"
        case NotFound:
            return "Not Found"
        case InternalServerError:
            return "Internal Server Error"
        default:
            panic(fmt.Sprintf("HTTP Status Code %d not supported", statusCode))
    }
}
