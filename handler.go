package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

func fetchData(id string) (string, error) {

	// this is a public api to mimic a concurrent task
	baseUrl := "https://jsonplaceholder.typicode.com/todos"
	if id != "" {
		baseUrl = "https://jsonplaceholder.typicode.com/todos/" + id
	}

	response, err := http.Get(baseUrl)
	if err != nil {
		return "", fmt.Errorf("error making external request: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {

	// mimic db persistence - eventually this will be a write to db
	log.Printf("mock persistence: %s", r.URL.Path)

	// upgrade to web socket protocol
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading to web socket protocol: %v", err)
		return
	}
	defer ws.Close()

	// increment counter
	atomic.AddInt64(&s.activeConnections, 1)
	defer atomic.AddInt64(&s.activeConnections, -1)

	// log connection count
	log.Printf("Active connections: %d", atomic.LoadInt64(&s.activeConnections))

	// simulation of second service
	go func() {

		// call public api
		response, err := fetchData(strconv.Itoa(rand.Intn(200) + 1))
		if err != nil {
			log.Printf("Write error: %v", err)
		}

		err = ws.WriteMessage(websocket.TextMessage, []byte("Transaction processed: "+response))
		if err != nil {
			log.Printf("Write error: %v", err)
		}
	}()

	// keep connection alive until process is done
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
