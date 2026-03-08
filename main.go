package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {

}

func main() {

	app := application.New(application.Options{
		Name:        "YAMP",
		Description: "Yet Another Music Player",
		Services:    []application.Service{},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

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

	go func() {

	}()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
