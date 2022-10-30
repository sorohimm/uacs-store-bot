package main

import (
	"context"

	"uacs_store_bot/internal/service/bot"
)

var version, buildTime string

func main() {
	app := bot.NewService()
	app.Init(context.Background(), "uacs-store-bot", version, buildTime)
}
