package main

import "log"

const APP_NAME = "QBT-Proxy"
const ENV_PREFIX = "QBP"

func main() {
	var cfg Config
	if err := cfg.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	var server Server
	server.Run(cfg)
}
