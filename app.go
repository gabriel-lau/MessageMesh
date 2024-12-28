package main

import (
	backend "MessageMesh/backend"
	models "MessageMesh/backend/models"
	"context"
	"fmt"
	"time"

	"github.com/wailsapp/wails"
)

const (
	blue   = "\033[34m"
	purple = "\033[35m"
	pink   = "\033[95m"
)

// App struct
type App struct {
	ctx     context.Context
	network backend.Network
	runtime *wails.Runtime
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.network.ConnectToNetwork()
}

func (a *App) WailsInit(runtime *wails.Runtime) error {
	go func() {
		for {
			select {
			case msg := <-a.network.ChatRoom.Inbound:
				runtime.Events.Emit("getMessage", msg.Message)
				fmt.Printf(blue+"[server.go]"+" ["+time.Now().Format("15:04:05")+"]"+reset+" Message: %s\n", msg.Message)
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// CHATCOMPONET
func (a *App) SendMessage(message string) {
	newMessage := models.Message{
		Sender:    "test",
		Receiver:  "test",
		Message:   message,
		Timestamp: "test",
	}

	a.network.SendMessage(newMessage.Message)
}
