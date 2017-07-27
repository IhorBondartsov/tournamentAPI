package main

import (
	"github.com/IhorBondartsov/tournamentAPI/src/server"
)

func main() {
	close := make(chan struct{})
	server := server.NewHTTPServer(port, host, close)
	server.DB = db

	server.Run()
	<-close
}
