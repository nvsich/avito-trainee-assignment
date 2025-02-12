package main

import "avito-shop/internal/app"

// TODO: think about adding as flag

const configPath = "./config/config.yaml"

func main() {
	app.Run(configPath)
}
