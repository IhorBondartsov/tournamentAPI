package main

import (
	"github.com/KharkivGophers/TeamplayerAPI/src/dao"
	"github.com/KharkivGophers/TeamplayerAPI/src/dao/DB"
)

var (
	host string = "0.0.0.0"
	port int32  = 8100

	//db dao.DAOInterface = &DB.SQL{}
	db dao.DAOInterface = &DB.StubDB{}
)
