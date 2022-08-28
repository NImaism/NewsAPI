package main

import (
	"Newism/internal/database"
	"Newism/internal/handler"
	"Newism/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	// INIT SERVER
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}

	server := gin.Default()

	mongodb := database.New()

	hnd := handler.New(logger.Named("handler.App"), store.New(mongodb))

	hnd.SetRouting(server)

	// RUN SERVER
	if err := server.Run(":" + os.Args[1]); err != nil {
		logger.Error("Can't run server", zap.Error(err))
	}
}
