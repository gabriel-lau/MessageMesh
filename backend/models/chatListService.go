package backend

type ChatListService struct {
	ChatServices []ChatService `json:"chatservices"`
}

func NewChatListService() *ChatListService {
	return &ChatListService{
		ChatServices: make([]ChatService, 0),
	}
}
