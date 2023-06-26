package app

import (
	"chat-server/internal/handlers"
	"chat-server/internal/repository"
	"chat-server/internal/route"
	"chat-server/internal/service"
	"chat-server/pkg/db"
	"chat-server/pkg/logs"
	"chat-server/pkg/server"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
)

func Run() {
	// yaml
	if err := initConfig(); err != nil {
		log.Fatalf("Error occured while init viper configs: %s", err.Error())
		return
	}
	var config db.Config
	if err := viper.UnmarshalKey("db", &config); err != nil {
		log.Fatalf("Error unmarshaling configs: %s", err)
	}
	config.Password = os.Getenv("DB_PASSWORD")

	db.StartDbConnection(config)
	defer db.CloseDbConnection()

	// logger
	logs.InitLogger()
	defer logs.CloseLoggerFile()

	userRep := repository.NewUserRepository(db.GetDBConn())
	roomRep := repository.NewRoomRepository(db.GetDBConn())
	memberRep := repository.NewMemberRepository(db.GetDBConn())
	msgRep := repository.NewMessageRepository(db.GetDBConn())
	userSvc := service.NewUserService(userRep)
	roomSvc := service.NewRoomService(roomRep, memberRep)
	authSvc := service.NewAuthService(userSvc)
	messageSvc := service.NewMessageService(msgRep)

	authHandler := handlers.NewAuthHandler(userSvc, authSvc, logs.GetLogger())
	roomHandler := handlers.NewRoomHandler(roomSvc, logs.GetLogger())
	chatHandler := handlers.NewChatHandler(messageSvc, logs.GetLogger())
	router := route.NewRouter(authHandler, roomHandler, chatHandler)
	httpPort := viper.GetString("port")

	srv := new(server.Server)
	go func() {
		if err := srv.Run(httpPort, router.SetupRouter()); err != nil {
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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
