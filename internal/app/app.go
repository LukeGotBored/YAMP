package app

import (
	"changeme/internal/appevents"
	"changeme/internal/application/services"
	"changeme/internal/db"
	"changeme/internal/infrastructure/persistence/sqlite"
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func getDBPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	appDir := filepath.Join(configDir, "YAMP")
	os.MkdirAll(appDir, 0755)
	return filepath.Join(appDir, "yamp.db")
}

func NewApp(assets embed.FS) (*Context, error) {
	dbService, err := db.NewService(getDBPath())
	if err != nil {
		return nil, err
	}

	// DI: Wiring up infrastructure and repositories
	trackRepo := sqlite.NewTrackRepository(dbService.DB, dbService.Queries)
	albumRepo := sqlite.NewAlbumRepository(dbService.Queries)
	artistRepo := sqlite.NewArtistRepository(dbService.Queries)

	// Application Use Cases
	libraryService := services.NewLibraryService(trackRepo, albumRepo, artistRepo)

	ctx := &Context{}
	opts := application.Options{
		Name:        AppName,
		Description: AppDescription,
		PanicHandler: func(details *application.PanicDetails) {
			log.Printf("Recovered backend panic: %v", details.Error)
			appevents.EmitRecoveredPanic(details)
		},
		Services: []application.Service{
			application.NewService(ctx),
			application.NewService(libraryService),
			application.NewService(dbService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	}
	ctx.app = application.New(opts)

	// Setup application context menus
	SetupContextMenus(ctx.app)

	ctx.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: AppName,
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
	if err := ctx.app.Run(); err != nil {
		return nil, err
	}
	return ctx, nil
}
