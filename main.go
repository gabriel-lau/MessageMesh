package main

import (
	appservices "MessageMesh/backend/models"

	"context"
	"embed"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS
var headless = false

const (
	yellow = "\033[33m"
	reset  = "\033[0m"
)

func main() {
	if headless {
		fmt.Println(yellow + "[main.go]" + reset + " Running in headless mode")
		app := NewApp()
		ctx := context.Background()
		app.startup(ctx)

		for {
			// app.SendMessage("Hello")
			time.Sleep(3 * time.Second)
		}
	}

	fmt.Println(yellow + "[main.go]" + reset + " Running in GUI mode")

	// Create an instance of the app structure
	app := NewApp()
	appService := *appservices.NewAppService()
	app.appService = appService

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "MessageMesh",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			&appService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
