package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/PCoelho07/gochat/protocol"
)

type Client struct {
	id   string
	name string
}

func Start(addr string) error {
	conn, err := connect(addr)
	if err != nil {
		return fmt.Errorf("start: %v", err)
	}

	defer conn.Close()

	input := bufio.NewReader(os.Stdin)
	fmt.Print("name: ")
	username, err := input.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading username: %v", err)
	}

    fmt.Printf(username)

	done := make(chan struct{})
	out := make(chan string)

	go Present(out)
	go readFromServer(conn, out, done)

    writeToServer(done, out, input, username, conn)

    return nil
}

func connect(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("connect: %v", err)
	}

	return conn, nil
}

func readFromServer(conn net.Conn, res chan string, done chan struct{}) {
	reader := bufio.NewReader(conn)
	var m protocol.Message

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			res <- fmt.Sprintf("readFromServer: %v", err)
			close(done)
			return
		}

		err = json.Unmarshal([]byte(msg), &m)
		if err != nil {
			res <- fmt.Sprintf("readFromServer: %v", err)
			close(done)
			return
		}

		res <- fmt.Sprintf("%s says: %s", m.SenderName, m.Text)
		res <- "> "
	}
}

func writeToServer(done <-chan struct{}, out chan<- string, stdin *bufio.Reader, username string, conn net.Conn) error { 
    for {
        select {
        case <-done:
            close(out)
            return nil
        default:
        }

        out <- "> "
        text, err := stdin.ReadString('\n')
        if err != nil {
            close(out)
            return fmt.Errorf("start: reading stdin: %w", err)
        }

        msg := protocol.Message{
            Text:   text,
            SenderName: strings.TrimSpace(username),
        }

        encoded, err := json.Marshal(msg)
        encoded = append(encoded, '\n')

        _, err = conn.Write(encoded)
        if err != nil {
            close(out)
            return fmt.Errorf("start: reading stdin: %w", err)
        }
    }
}
