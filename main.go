package main

import (
	debug "MessageMesh/debug"
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	if debug.IsHeadless {
		debug.Log("main", "Running in headless mode")
		app := NewApp()
		ctx := context.Background()
		app.startup(ctx)
		// Keep the app running
		select {}
	}

	debug.Log("main", "Running in GUI mode")

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "MessageMesh",
		Width:  1080,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		debug.Log("error", err.Error())
	}
}
