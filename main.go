package main

import (
	"authentication/config"
	"authentication/routes"
)

func main() {
	cfg := config.GetConfig()
	router := routes.GetRouter()
	router.Run(cfg.ServerAddr)
}
