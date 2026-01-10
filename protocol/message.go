package protocol

type Message struct {
	Text       string `json:"text"`
	Sender     string `json:"sender"`
	SenderName string `json:"sender_name"`
}
