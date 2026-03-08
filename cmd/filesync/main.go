package main

import (
	"log"
	"os"

	"github.com/lonelyday/filesync/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}
