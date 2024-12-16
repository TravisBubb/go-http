package http

import (
	"encoding/json"
	"errors"
	"strings"
)

type Context struct {
	Request  *HttpRequest
	Response *HttpResponse
	handlers HandlerPipeline
}

func createContext(r *HttpRequest) Context {
	return Context{
		Request:  r,
		Response: nil,
		handlers: make(HandlerPipeline, 0),
	}
}

func (c *Context) addHandler(handler Handler) {
	c.handlers = append(c.handlers, handler)
}

func (c *Context) BindRequest(object any) error {
	if c.Request == nil {
		return errors.New("invalid request")
	}

	decoder := json.NewDecoder(strings.NewReader(c.Request.Content))

	return decoder.Decode(object)
}

func (c *Context) Ok(content string) {
	c.Response = &HttpResponse{
		StatusCode: OK,
		Content:    content,
		Headers:    make(map[string]string),
		Protocol:   c.Request.Protocol,
	}
}

func (c *Context) BadRequest(content string) {
	c.Response = &HttpResponse{
		StatusCode: BadRequest,
		Content:    content,
		Headers:    make(map[string]string),
		Protocol:   c.Request.Protocol,
	}
}
