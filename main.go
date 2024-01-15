package main

import (
	"context"
	"fmt"
	"path/filepath"
	"pidorator-bot/bot"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	absPath, err := filepath.Abs(".env")
	if err != nil {
		fmt.Printf("Error on path to .env creation, %s\n", err)
		return
	}
	if err := godotenv.Load(absPath); err != nil {
		fmt.Printf("Cannot read .env: %s\n", err)
	}

	b, err := bot.New(ctx)
	if err != nil {
		fmt.Printf("Error. Can't create client: %s", err)
		return
	}

	b.InitHandlers()
	b.Start(ctx)
}
