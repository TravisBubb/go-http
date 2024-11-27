package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Represents an instance of a server
type server struct {
	host     string
	port     uint16
	listener net.Listener
}

// Creates a new TCP server that listens on the provided host and port number
func CreateServer(host string, port uint16) *server {
	return &server{
		host: host,
		port: port,
	}
}

// Starts the server
func (s *server) Start() error {
	address := fmt.Sprintf("%s:%d", s.host, s.port)
	var err error

	s.listener, err = net.Listen("tcp4", address)
	if err != nil {
		log.Printf("Failed to bind server to address %s", address)
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	go s.serve(ctx)

	// Handle os signals to stop the server
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChannel
		cancel()
		s.stop()
	}()

	<-ctx.Done()

	return err
}

// Stops the server
func (s *server) stop() {
	s.listener.Close()
	fmt.Println("tcp listener closed.")
}

// Executes the main server loop that accepts new incoming connections
func (s *server) serve(ctx context.Context) {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			log.Println("Error accepting connection:", err)
		} else {
			go func() {
				s.handleConnection(ctx, connection)
			}()
		}
	}
}

// Handles any connections made with the server
func (s *server) handleConnection(ctx context.Context, connection net.Conn) {
	defer connection.Close()

	request, err := GetRequestFromConnection(ctx, connection)
	if err != nil {
		log.Println("Error occurred attempting to parse request:", err)
		return
	}

	log.Println("Request:", request)
}
