package main

import (
	"github.com/KharkivGophers/TeamplayerAPI/src/server"
)

func main() {
	close := make(chan struct{})
	server := server.NewHTTPServer(port, host)
	server.DB = db

	server.Run()
	<-close
}
