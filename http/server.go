package http

import (
    "context"
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type server struct {
    host      string
    port      uint16
    listener  net.Listener
}

func CreateServer(host string, port uint16) *server {
    return &server{
        host: host,
        port: port,
    }
}

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

func (s *server) stop() {
    s.listener.Close()
    fmt.Println("tcp listener closed.")
}

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
                s.handleConnection(connection, ctx)
            }()
        }
    }
}

func (s *server) handleConnection(connection net.Conn, ctx context.Context) {
    defer connection.Close()
    buffer := make([]byte, 2048)

ReadLoop:
    for {
        if ctx.Err() != nil {
            return
        }

        err := connection.SetDeadline(time.Now().Add(200 * time.Millisecond))
        if err != nil {
            log.Println("Unexpected error attempting to set connection timeout", err)
            return
        }

        n, err := connection.Read(buffer)
        if err != nil {
            if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                continue ReadLoop
            } else if err != io.EOF {
                log.Println("Unexpected error reading from connection", err)
                return
            }
        }

        if n == 0 {
            return
        }

        log.Printf("Received from %v: %s", connection.RemoteAddr(), string(buffer[:n]))
    }
}
