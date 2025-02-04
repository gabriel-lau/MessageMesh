package main

import (
	backend "MessageMesh/backend"
	debug "MessageMesh/debug"
	"context"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	network backend.Network
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

	// Events Emitter
	go func() {
		debug.Log("app", "Wails events emitter started")
		refreshticker := time.NewTicker(time.Second)
		defer refreshticker.Stop()

		if GetEnvVar("HEADELESS") == "false" {
			for {
				select {
				case msg := <-a.network.ChatRoom.Inbound:
					runtime.EventsEmit(a.ctx, "getMessage", msg.Message)
					// debug.Log("app", "Message: "+msg.Message)
					time.Sleep(1 * time.Second)
				case <-refreshticker.C:
					runtime.EventsEmit(a.ctx, "getPeersList", a.network.ChatRoom.PeerList())
				}
			}
		}
	}()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// CHATCOMPONET
func (a *App) SendMessage(message string) {
	a.network.SendMessage(message)
}
