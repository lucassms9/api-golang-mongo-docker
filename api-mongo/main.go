package main

import (
	"api-mongo/server"
)

func main() {
	var handle = server.Init()
	go handle.StartServer()
}
