package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	activeConnections int64
	upgrader          websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func main() {
	server := NewServer()

	// API endpoint that handles initial request to process
	http.HandleFunc("/process", server.handleRequest)

	// Load test endpoint - this will begin the creation of 1000 concurrent requests
	http.HandleFunc("/loadtest", func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		startTime := time.Now()

		// Simulate 1000 concurrent requests
		for i := 0; i < 1000; i++ {
			// add iteration to wait group
			wg.Add(1)
			go func(requestNum int) {
				defer wg.Done()

				// Create WebSocket connection
				url := "ws://localhost:8080/process"
				c, _, err := websocket.DefaultDialer.Dial(url, nil)
				if err != nil {
					log.Printf("Dial error: %v", err)
					return
				}
				defer c.Close()

				// Wait for response
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Printf("Read error: %v", err)
					return
				}
				log.Printf("Request %d completed: %s", requestNum, message)
			}(i)
		}

		wg.Wait()
		duration := time.Since(startTime)
		fmt.Fprintf(w, "Load test completed in %v", duration)
	})

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
