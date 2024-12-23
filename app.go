package main

import (
	"MessageMesh/backend"
	"context"
	"fmt"
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
	// runtime.EventsEmit(a.ctx, "emitMyEvent", a.network.ChatRoom.pstopic.ListPeers())
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SendMessage(message string) {
	a.network.SendMessage(message)
}
