package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	"chat-server/config"
	"chat-server/internal/handlers"
	"chat-server/internal/repository"
	"chat-server/internal/route"
	"chat-server/internal/service"
	"chat-server/pkg/db"
	"chat-server/pkg/logger"
	"chat-server/pkg/redis"
	"chat-server/pkg/server"
)

func Run() {
	// yaml
	cfg := config.Config{}
	if err := cfg.Parse(); err != nil {
		log.Fatalf("Error while parsing yml file: %v", err)
		return
	}
	// db
	db.StartDbConnection(cfg.GetDBConfig())
	defer db.CloseDbConnection()

	// redis
	redis.StartRedisConnection(cfg.GetRedisConfig())
	defer redis.CloseRedisConnection()

	// logger
	logger.InitLogger()
	defer logger.CloseLoggerFile()

	userRep := repository.NewUserRepository(db.GetDBConn())
	roomRep := repository.NewRoomRepository(db.GetDBConn())
	memberRep := repository.NewMemberRepository(db.GetDBConn())
	msgRep := repository.NewMessageRepository(db.GetDBConn())
	tokenCacheRep := repository.NewTokenCacheRepository(redis.GetRedisConn())
	userCacheRep := repository.NewUserCacheRepository(redis.GetRedisConn())

	userSvc := service.NewUserService(userRep, userCacheRep)
	roomSvc := service.NewRoomService(roomRep, memberRep)
	tokenSvc := service.NewTokenService(cfg.GetTSConfig(), tokenCacheRep)
	authSvc := service.NewAuthService(userSvc, tokenSvc)
	messageSvc := service.NewMessageService(msgRep)

	authHandler := handlers.NewAuthHandler(userSvc, authSvc, tokenSvc, logger.GetLogger())
	roomHandler := handlers.NewRoomHandler(roomSvc, logger.GetLogger())
	chatHandler := handlers.NewChatHandler(messageSvc, logger.GetLogger())
	router := route.NewRouter(authHandler, roomHandler, chatHandler)
	httpPort := cfg.GetServerPort()

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
