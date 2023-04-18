package main

import (
	"embed"

	"github.com/joho/godotenv"
	"github.com/templari/shire-client/util"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "shire-client",
		Width:            1024,
		MinWidth:         434,
		Height:           768,
		MinHeight:        318,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,

		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameAqua,
			WindowIsTranslucent:  true,
			WebviewIsTransparent: true,
		},
	})

	if err := godotenv.Load(); err != nil {
		util.Logger.Fatal("Error loading .env file")
	}

	if err != nil {
		println("Error:", err.Error())
	}
}
