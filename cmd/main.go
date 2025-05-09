package main

import (
	"github.com/Yorshik/Go-Final/internal/agent"
	"github.com/Yorshik/Go-Final/internal/server"
)

func main() {
	go agent.StartServer()
	server.StartServer()
}
