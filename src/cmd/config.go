package main

import (
	"github.com/KharkivGophers/TeamplayerAPI/src/dao"
	"github.com/KharkivGophers/TeamplayerAPI/src/dao/DB"
)

var (
	host string = "0.0.0.0"
	port int32  = 8100


	baseUser string = "root"
	basePassword string
	baseURI string
	baseTypeConn string = "tcp"

	//db dao.DAOInterface = &DB.SQL{ Password: basePassword, URI:baseURI, User: baseUser, TypeConn: baseTypeConn}
	db dao.DAOInterface = &DB.StubDB{}
)
