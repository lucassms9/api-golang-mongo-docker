package main

import (
	"api-mongo/server"
)

func main() {
	var handle = server.Init()
	handle.StartServer()
}
