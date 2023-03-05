package main

import (
	"im-instance/router"
)

func main() {
	r := router.Router()

	if err := r.Run(":7898"); err != nil {
		panic(err.Error())
	}
}
