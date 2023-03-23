package main

import (
	"client/common"
	"client/router"
	"client/wsHandler"
	"log"
)

func main() {
	wsHandler.InitHandler()
	ginServer := router.InitRouter()
	log.Println("http://localhost:10000/login")
	//
	ginServer.Run(":" + common.Client_port)

}
