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
    handleConnection func(context.Context, net.Conn)
}

// Creates a new TCP server that listens on the provided host and port number
func CreateServer(host string, port uint16, handleConnection func(context.Context, net.Conn)) *server {
	return &server{
		host: host,
		port: port,
        handleConnection: handleConnection,
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
            go func(){
                defer connection.Close()
                s.handleConnection(ctx, connection)
            }()
		}
	}
}
