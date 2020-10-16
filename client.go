package main

import (
	"github.com/gorilla/websocket"
)

// Client Structure
type Client struct {
	// Connection Object
	conn *websocket.Conn

	// Channel to send data
	send chan []byte

	// Reference to room instance
	room *Room
}

// Method to read data from the socket
func (c *Client) Read() {
	defer c.conn.Close()
	for {
		// Read Message from the socket connection
		_, msg, err := c.conn.ReadMessage()
		// If error, return
		if err != nil {
			return
		}
		// If no error, transmit to connectionPool of the room
		c.room.forwardQueue <- msg
	}
}

// Method to write data into the socket
func (c *Client) Write() {
	defer c.conn.Close()
	// Iterate over the channel of the client
	for msg := range c.send {
		// write the message to the socket
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		// Handle Error
		if err != nil {
			return
		}
	}
}
