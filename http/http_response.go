package http

// An HttpResponse is the representation of an HTTP request.
type HttpResponse struct {
    StatusCode HttpStatusCode
	Protocol string
	Headers  map[string]string
	Content  string
}
