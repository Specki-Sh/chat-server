package app

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"chat-server/config"
	"chat-server/internal/handlers"
	"chat-server/internal/repository"
	"chat-server/internal/route"
	"chat-server/internal/service"
)

func RouterFactory(
	logger *logrus.Logger,
	conn *sql.DB,
	redisClient *redis.Client,
	cfg config.Config,
) *route.Router {
	authHandler := authHandlerFactory(logger, conn, redisClient, cfg.GetTSConfig())
	roomHandler := roomHandlerFactory(logger, conn)
	chatHandler := chatHandlerFactory(logger, conn)
	return route.NewRouter(authHandler, roomHandler, chatHandler)
}

func authHandlerFactory(
	logger *logrus.Logger,
	conn *sql.DB,
	redisClient *redis.Client,
	tsConfig *service.TSConfig,
) *handlers.AuthHandler {
	userRep := repository.NewUserRepository(conn)
	userCacheRep := repository.NewUserCacheRepository(redisClient)
	tokenCacheRep := repository.NewTokenCacheRepository(redisClient)

	userSvc := service.NewUserService(userRep, userCacheRep)
	tokenSvc := service.NewTokenService(tsConfig, tokenCacheRep)
	authSvc := service.NewAuthService(userSvc, tokenSvc)

	return handlers.NewAuthHandler(userSvc, authSvc, tokenSvc, logger)
}

func roomHandlerFactory(logger *logrus.Logger, conn *sql.DB) *handlers.RoomHandler {
	roomRep := repository.NewRoomRepository(conn)
	memberRep := repository.NewMemberRepository(conn)
	roomSvc := service.NewRoomService(roomRep, memberRep)
	return handlers.NewRoomHandler(roomSvc, logger)
}

func chatHandlerFactory(logger *logrus.Logger, conn *sql.DB) *handlers.ChatHandler {
	msgRep := repository.NewMessageRepository(conn)
	messageSvc := service.NewMessageService(msgRep)
	return handlers.NewChatHandler(messageSvc, logger)
}
