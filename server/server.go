package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/PCoelho07/gochat/protocol"
	"github.com/google/uuid"
)

type Server struct {
	Clist []Client
}

type Client struct {
	id   string
	name string
	conn net.Conn
	send chan []byte
}

func Start(port string) error {
	slog.Info("starting server", "port", port)

	ln, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("start: %v", err)
	}

	hub := NewHub()
	go hub.Run()

	for {
		c, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("start: %v", err)
		}

		cl := &Client{
			id:   uuid.NewString(),
			conn: c,
			send: make(chan []byte, 256),
		}

		slog.Info("client connected", "id", cl.id)
		hub.register <- cl

		go handleRead(cl, hub)
		go handleWrite(cl)
	}
}

func handleRead(c *Client, h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()

	reader := bufio.NewReader(c.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		var m protocol.Message
		if err := json.Unmarshal([]byte(msg), &m); err != nil {
			continue
		}

        c.name = strings.TrimSpace(m.SenderName)
		m.Sender = c.id
		h.broadcast <- m
	}
}

func handleWrite(client *Client) {
	for msg := range client.send {
		_, err := client.conn.Write(msg)
		if err != nil {
			return
		}
	}
}
