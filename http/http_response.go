package http

// An HttpResponse is the representation of an HTTP request.
type HttpResponse struct {
    StatusCode HttpStatusCode
	Protocol string
	Headers  map[string]string
	Content  string
}

// Creates a new HttpResponse
func CreateResponse(statusCode HttpStatusCode, content string, headers map[string]string, protocol string) HttpResponse {
    response := HttpResponse {
        StatusCode: statusCode,
        Content: content,
        Headers: headers,
        Protocol: protocol,
    }

    return response
}
