package backend

type AppService struct {
	ChatListService ChatListService `json:"chatListService"`
}

func NewAppService() *AppService {
	return &AppService{
		ChatListService: *NewChatListService(),
	}
}
