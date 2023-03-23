package main

import (
	"fmt"
	"server/common"
	"server/wsHandler"
)

func main() {
	db := common.InitDB()
	defer db.Close()
	wsHandler.InitHandler()

	fmt.Println("startListen")

}
