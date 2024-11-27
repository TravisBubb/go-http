package http

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// An HttpRequest is the representation of an HTTP request.
type HttpRequest struct {
	Method   HttpMethod
	Path     string
	Protocol string
	Headers  map[string]string
	Content  string
}

// Creates a new HttpRequest from a given TCP connection.
func GetRequestFromConnection(ctx context.Context, connection net.Conn) (HttpRequest, error) {
	request := HttpRequest{}
	if ctx.Err() != nil {
		return HttpRequest{}, ctx.Err()
	}

	err := connection.SetDeadline(time.Now().Add(1000 * time.Millisecond))
	if err != nil {
		log.Println("Unexpected error attempting to set connection timeout", err)
		return HttpRequest{}, err
	}

	reader := bufio.NewReader(connection)

	err = request.readHttpMetadata(reader)
	if err != nil {
		log.Println("Error occurred while reading HTTP metadata:", err)
		return HttpRequest{}, nil
	}

	err = request.readHeaders(reader)
	if err != nil {
		log.Println("Error occurred while reading HTTP headers:", err)
		return HttpRequest{}, err
	}

	err = request.readContent(reader)
	if err != nil {
		log.Println("Error occurred while reading HTTP request content:", err)
		return HttpRequest{}, nil
	}

	return request, nil
}

// Load the HTTP request information using the provided bufio.Reader
func (request *HttpRequest) readHttpMetadata(reader *bufio.Reader) error {
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Error reading HTTP metadata: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) != 3 {
		return fmt.Errorf("Invalid HTTP metadata: %w", err)
	}

	switch parts[0] {
	case "GET":
		request.Method = GET
	case "POST":
		request.Method = POST
	case "PUT":
		request.Method = PUT
	case "DELETE":
		request.Method = DELETE
	default:
		return fmt.Errorf("HTTP method [%s] is not currently supported", parts[0])
	}

	request.Path = parts[1]
	request.Protocol = parts[2]
	return nil
}

// Load the HTTP headers using the provided bufio.Reader
func (request *HttpRequest) readHeaders(reader *bufio.Reader) error {
	request.Headers = make(map[string]string)
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading headers:", err)
			return err
		}

		headerLine = strings.TrimSpace(headerLine)
		if headerLine == "" {
			break // End of header section
		}

		headerParts := strings.SplitN(headerLine, ":", 2)
		if len(headerParts) != 2 {
			log.Println("Invalid header:", headerLine)
			return fmt.Errorf("Invalid header: %s", headerLine)
		}

		request.Headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
	}
	return nil
}

// Load the HTTP request content using the provided bufio.Reader
func (request *HttpRequest) readContent(reader *bufio.Reader) error {

	val, ok := request.Headers["Content-Length"]
	if !ok {
		return nil
	}

	contentLength, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("Invalid Content-Length header value: %s", val)
	}

	buffer := make([]byte, contentLength)
	for i := 0; i < contentLength; i++ {
		byte, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("Error trying to parse request: %w", err)
		}
		buffer[i] = byte
	}

	request.Content = string(buffer)
	return nil
}
