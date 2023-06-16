package app

import (
	"chat-server/internal/handlers"
	"chat-server/internal/repository"
	"chat-server/internal/service"
	"chat-server/pkg/db"
	"chat-server/pkg/server"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	config := db.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "password",
		DBName:   "mydb",
		SSLMode:  "disable",
	}
	db.StartDbConnection(config)
	defer db.CloseDbConnection()

	userRep := repository.NewUserRepository(db.GetDBConn())
	userSvc := service.NewUserService(userRep)
	authSvc := service.NewAuthService(userSvc)

	authHandler := handlers.NewAuthHandler(userSvc, authSvc)

	httpPort := "8000"

	srv := server.NewServer(authHandler)
	go func() {
		if err := srv.Run(httpPort); err != nil {
			log.Fatalf("Error occured while running http server: %s", err.Error())
			return
		}
	}()

	fmt.Println("App Started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	fmt.Println("Shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occurred on server shutting down: %s", err.Error())
	}
}
