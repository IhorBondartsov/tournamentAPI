package main

import (
	"github.com/IhorBondartsov/tournamentAPI/src/dao"
	"github.com/IhorBondartsov/tournamentAPI/src/dao/DB"
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
