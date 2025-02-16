package main

import (
	"avito-shop/internal/app"
	"flag"
)

func main() {
	var envPath string
	flag.StringVar(&envPath, "env-path", ".env", "path to .env")
	flag.Parse()

	app.Run(envPath)
}
