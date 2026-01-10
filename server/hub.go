package server

import (
	"encoding/json"
	"log/slog"

	"github.com/PCoelho07/gochat/protocol"
)

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan protocol.Message
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan protocol.Message),
	}
}

func (h *Hub) Run() error {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		case msg := <-h.broadcast:
            out, err := json.Marshal(msg)
            if err != nil { 
                slog.Error("marshal message", "error", err)
                continue
            }

            out = append(out, '\n')

			for _, client := range h.clients {
                if msg.Sender == client.id {
                    continue
                }

				select {
                case client.send <- out:
				default:
					close(client.send)
					delete(h.clients, client.id)
				}
			}
		}
	}
}
