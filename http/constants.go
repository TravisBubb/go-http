package http

// Enum for HTTP methods
type HttpMethod int

const (
	GET HttpMethod = iota
	POST
	PUT
	DELETE
)

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
