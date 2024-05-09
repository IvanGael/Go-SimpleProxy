package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <listen_addr> <target_addr>\n", os.Args[0])
		os.Exit(1)
	}

	listenAddr := os.Args[1]
	targetAddr := os.Args[2]

	// Listen for incoming connections
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening: %s\n", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("Proxy server listening on %s\n", listenAddr)

	for {
		// Accept an incoming connection
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %s\n", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleClientRequest(clientConn, targetAddr)
	}
}

func handleClientRequest(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	// Connect to the target server
	serverConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to target: %s\n", err)
		return
	}
	defer serverConn.Close()

	// Forward data between client and target server
	go func() {
		_, err := io.Copy(serverConn, clientConn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error copying to server: %s\n", err)
		}
	}()

	_, err = io.Copy(clientConn, serverConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying to client: %s\n", err)
	}
}

//go run main.go :8080 example.com:80
