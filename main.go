package main

import (
	"im-instance/router"
	"im-instance/utils"
)

func main() {
	utils.InitMysql()

	r := router.Router()

	if err := r.Run(":7898"); err != nil {
		panic(err.Error())
	}
}
