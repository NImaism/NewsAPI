package main

import (
	database "Newism/Database"
	"Newism/Middleware"
	"Newism/Routing"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	// INIT SERVER
	Server := gin.Default()
	database.Connect()

	//// SET LIMIT
	Server.Use(Middleware.RateLimit(1, 1))

	// SET ROUTING
	Routing.SetRouting(Server)
	// LOAD FILE

	if err := Server.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}
