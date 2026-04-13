package app

import "github.com/wailsapp/wails/v3/pkg/application"

type Context struct {
	app *application.App
}

// ShowLibraryItemMenu opens the native context menu for a library track.
func (c *Context) ShowLibraryItemMenu(x, y int, itemId string) {
	if c.app == nil {
		return
	}

	// Get the frontmost window
	if w := c.app.Window.Current(); w != nil {
		w.OpenContextMenu(&application.ContextMenuData{
			Id:   "library_item_menu",
			X:    x,
			Y:    y,
			Data: itemId,
		})
	}
}
