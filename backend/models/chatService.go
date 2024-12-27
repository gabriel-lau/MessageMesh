package backend

type ChatService struct {
	Recipient string    `json:"recipient"`
	Meessages []Message `json:"messages"`
}

type Message struct {
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	TimeStamp string `json:"timestamp"`
}

func (cs *ChatService) SendMessage(message string) {
	cs.Meessages = append(cs.Meessages, Message{
		Recipient: cs.Recipient,
		Sender:    "Me",
		Message:   message,
		TimeStamp: "Now",
	})
}

func NewChatService() *ChatService {
	return &ChatService{
		Meessages: make([]Message, 0),
	}
}
