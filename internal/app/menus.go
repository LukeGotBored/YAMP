package app

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// SetupContextMenus creates the context menus available for the frontend to trigger.
func SetupContextMenus(app *application.App) {
	// Create a new context menu
	libraryMenu := app.ContextMenu.New()

	libraryMenu.Add("Play").OnClick(func(ctx *application.Context) {
		log.Println("Play clicked")
		app.Event.Emit("play_triggered", "some_data")
	})
	libraryMenu.Add("Add to Queue").OnClick(func(ctx *application.Context) {
		log.Println("Add to Queue clicked")
	})
	libraryMenu.AddSeparator()
	libraryMenu.Add("Go to Artist").OnClick(func(ctx *application.Context) {
		log.Println("Go to Artist clicked")
	})

	// Register it with the application under an ID
	app.ContextMenu.Add("library_item_menu", libraryMenu)
}