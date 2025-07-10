package api_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/c2FmZQ/quic-api"
)

// myListener is a trivial wrapper around [api.Listener] that increases a
// counter when data is read from a stream on a received connection.
type myListener struct {
	api.Listener
	counter *atomic.Int32
}

func (ln *myListener) Accept(ctx context.Context) (api.Conn, error) {
	conn, err := ln.Listener.Accept(ctx)
	if err != nil {
		return nil, err
	}
	return &myConn{Conn: conn, counter: ln.counter}, nil
}

// myConn is a trivial wrapper around [api.Conn] that increases a
// counter when data is read from a stream on a received connection.
type myConn struct {
	api.Conn
	counter *atomic.Int32
}

func (c *myConn) AcceptStream(ctx context.Context) (api.Stream, error) {
	stream, err := c.Conn.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	return &myStream{Stream: stream, counter: c.counter}, nil
}

func (c *myConn) OpenStream() (api.Stream, error) {
	stream, err := c.Conn.OpenStream()
	if err != nil {
		return nil, err
	}
	return &myStream{Stream: stream, counter: c.counter}, nil
}

// myStream is a trivial wrapper around [api.Stream] that increases a
// counter when data is read from a stream on a received connection.
type myStream struct {
	api.Stream
	counter *atomic.Int32
}

func (s *myStream) Read(b []byte) (int, error) {
	n, err := s.Stream.Read(b)
	s.counter.Add(int32(n)) // Count the number of bytes read
	return n, err
}

func Example() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverReadCounter := new(atomic.Int32)
	clientReadCounter := new(atomic.Int32)

	cert, err := newCert("example.com")
	if err != nil {
		log.Fatalf("newCert: %v", err)
	}

	// Start the server.
	ln, err := api.ListenAddr("localhost:0", &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"foo"},
	}, nil)
	if err != nil {
		log.Fatalf("ListenAddr: %v", err)
	}
	// Here, we replace ln with our own api.Listener implementation
	// which simply increments the counter for every byte read from
	// a stream.
	ln = &myListener{Listener: ln, counter: serverReadCounter}
	go quicServer(ctx, ln)

	pool := x509.NewCertPool()
	pool.AddCert(cert.Leaf)

	// Open a QUIC connection to the server.
	conn, err := api.DialAddr(ctx, ln.Addr().String(), &tls.Config{
		ServerName: "example.com",
		RootCAs:    pool,
		NextProtos: []string{"foo"},
	}, nil)
	if err != nil {
		log.Fatalf("DialAddr: %v", err)
	}
	// We also count the number of bytes read by the client.
	conn = &myConn{Conn: conn, counter: clientReadCounter}

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendMessage(conn, "Hello World!\n")
		}()
	}
	wg.Wait()

	fmt.Printf("Server Read Counter: %d\n", serverReadCounter.Load())
	fmt.Printf("Client Read Counter: %d\n", clientReadCounter.Load())
	// Output:
	// Server Read Counter: 1300
	// Client Read Counter: 400
}

// sendMessage opens a bidirectional stream, sends one message, and reads the
// answers before closing the stream. It doesn't know anything about the myConn
// implementation.
//
// When stream.Read() is called, the underlying counter is incremented.
func sendMessage(conn api.Conn, message string) {
	stream, err := conn.OpenStream()
	if err != nil {
		log.Fatalf("OpenStream: %v", err)
	}
	defer stream.Close()

	if _, err := stream.Write([]byte(message)); err != nil {
		log.Fatalf("Client Write: %v", err)
	}
	log.Printf("Client sent %q (len=%d)", message, len(message))

	answer := make([]byte, 1024)
	n, err := stream.Read(answer)
	if err != nil && !errors.Is(err, io.EOF) {
		log.Fatalf("Client Read: %v", err)
	}
	log.Printf("Client read %q (len=%d)", answer[:n], n)
}

// quicServer represents an arbitrary QUIC server that uses an [api.Listener] as
// argument. It doesn't know anything about the actual myListener
// implementation.
//
// When stream.Read() is called, the underlying counter is incremented.
func quicServer(ctx context.Context, ln api.Listener) {
	for {
		conn, err := ln.Accept(ctx)
		if err != nil {
			log.Printf("Accept: %v", err)
			return
		}
		go func(ctx context.Context, conn api.Conn) {
			for {
				stream, err := conn.AcceptStream(ctx)
				if err != nil {
					log.Printf("AcceptStream: %v", err)
					return
				}
				go func(stream api.Stream) {
					buf := make([]byte, 1024)
					for {
						n, err := stream.Read(buf)
						if n > 0 {
							if _, err := stream.Write([]byte("ack\n")); err != nil {
								log.Printf("Server Write: %v", err)
								return
							}
						}
						if err != nil {
							if !errors.Is(err, io.EOF) {
								log.Printf("Server Read: %d, %v", n, err)
							}
							break
						}
					}
					if err := stream.Close(); err != nil {
						log.Printf("Server Close: %v", err)
					}
				}(stream)
			}
		}(ctx, conn)
	}
}
