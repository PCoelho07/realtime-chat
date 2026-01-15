package client

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
