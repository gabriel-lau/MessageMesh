package main

import (
	debug "MessageMesh/debug"
	"context"
	"embed"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	if GetEnvVar("HEADLESS") == "true" {
		debug.Log("main", "Running in headless mode")
		app := NewApp()
		ctx := context.Background()
		app.startup(ctx)
		for {
			app.SendMessage("Hello I am " + GetEnvVar("USERNAME"))
			time.Sleep(30 * time.Second)
		}
	}

	debug.Log("main", "Running in GUI mode")

	// Create an instance of the app structure
	app := NewApp()

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
		},
	})

	if err != nil {
		debug.Log("error", err.Error())
	}
}

func GetEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		debug.Log("error", "Error loading .env file")
	}
	return os.Getenv(key)
}
