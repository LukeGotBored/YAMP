package main

import (
	"embed"
	_ "embed"
	"log"

	"changeme/internal/appevents"
	"changeme/internal/debug"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	println("Starting YAMP...")
}

func main() {
	opts := application.Options{
		Name:        "YAMP",
		Description: "Yet Another Music Player",
		PanicHandler: func(details *application.PanicDetails) {
			log.Printf("Recovered backend panic: %v", details.Error)
			appevents.EmitRecoveredPanic(details)
		},
		Services: []application.Service{
			application.NewService(&debug.Service{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	}

	app := application.New(opts)

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "YAMP",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(25, 25, 25),
		URL:              "/",
		Width:            1060,
		Height:           680,
		MinWidth:         760,
		MinHeight:        500,
		Frameless:        true,
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
