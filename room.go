package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Constants for upgrader object
const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// Room Structure
type Room struct {
	// Channel to store all the messages to be forward
	forwardQueue chan []byte

	// Channel for clients wishing to join
	join chan *Client

	// Channel for clients wishing to leave
	leaves chan *Client

	// Map to hold all current clients in the room
	clientStore map[*Client]bool
}

// Upgrader instance to upgrade an HTTP Connection to a WebSocket Connection
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

// ServeHTTP Method for room type to convert into http.Handler interface
func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// Upgrade function to elevate HTTP connection to WebSocket
	sock, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Client Instance
	client := &Client{
		conn: sock,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}

	// Add to join channel
	r.join <- client
	// Add to leaves queue once, main thread exits
	defer func() { r.leaves <- client }()
	// Write then read concurrently
	go client.Write()
	client.Read()
}

// Helper function to run the room and handle clients
func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			{
				// A Client that wants to join
				r.clientStore[client] = true
			}
		case client := <-r.leaves:
			{
				// A client that wants to leave
				r.clientStore[client] = false
				// Delete the entry from the map
				delete(r.clientStore, client)
				// Close the send channel of the client
				close(client.send)
			}
		case msg := <-r.forwardQueue:
			{
				// Send the message to all the clients in the clientStore
				for client := range r.clientStore {
					client.send <- msg
				}
			}
		}

	}
}

// Helper function to create a new room
func newRoom() *Room {
	return &Room{
		forwardQueue: make(chan []byte),
		join:         make(chan *Client),
		leaves:       make(chan *Client),
		clientStore:  make(map[*Client]bool),
	}
}
