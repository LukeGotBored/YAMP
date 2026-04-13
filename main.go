package main

import (
	"changeme/internal/app"
	"embed"
	"log"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	println("Starting YAMP...")
}

func main() {
	ctx, err := app.NewApp(assets)
	if err != nil {
		log.Fatal(err)
	}
	_ = ctx
}
