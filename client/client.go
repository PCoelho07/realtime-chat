package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/PCoelho07/gochat/protocol"
)

type Client struct {
	id   string
	name string
}

type EventType int

const (
	EventInput EventType = iota
	EventServer
	EventQuit
)

type Event struct {
	Type EventType
	Data string
}

func Start(addr string) error {
	conn, err := connect(addr)
	if err != nil {
		return fmt.Errorf("start: %v", err)
	}

	defer conn.Close()

	events := make(chan Event)

	stdin := bufio.NewReader(os.Stdin)
	fmt.Print("name: ")
	username, err := stdin.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading username: %v", err)
	}

	fmt.Printf(username)

	go readFromServer(conn, events)
	go readFromUser(stdin, events)

	fmt.Print("> ")

	for {
		ev := <-events

		switch ev.Type {
		case EventInput:
			msg := protocol.Message{
				Text:       ev.Data,
				SenderName: strings.TrimSpace(username),
			}

			encoded, err := json.Marshal(msg)
			encoded = append(encoded, '\n')

			if _, err = conn.Write(encoded); err != nil {
				return fmt.Errorf("start: reading stdin: %w", err)
			}
			fmt.Print("> ")

		case EventServer:
			var msg protocol.Message

			if err := json.Unmarshal([]byte(ev.Data), &msg); err != nil {
				return fmt.Errorf("unmarshal from server: %v", err)
			}

			fmt.Printf("%s says: %s", msg.SenderName, msg.Text)
			fmt.Print("> ")
		case EventQuit:
			slog.Error(ev.Data)
			return fmt.Errorf("error: %v", ev.Data)
		}
	}
}

func connect(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("connect: %v", err)
	}

	return conn, nil
}

func readFromServer(conn net.Conn, events chan<- Event) {
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			events <- Event{Type: EventQuit, Data: err.Error()}
			return
		}

		events <- Event{Type: EventServer, Data: msg}
	}
}

func readFromUser(stdin *bufio.Reader, events chan<- Event) {
	for {
		msg, err := stdin.ReadString('\n')
		if err != nil {
			events <- Event{Type: EventQuit, Data: err.Error()}
			return
		}

		events <- Event{Type: EventInput, Data: msg}
	}
}
