package server

import (
	"chat-server/internal/handlers"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func NewServer(authHandler *handlers.AuthHandler, roomHandler *handlers.RoomHandler, chatHandler *handlers.ChatHandler) *Server {
	return &Server{
		httpServer:  &http.Server{},
		route:       gin.New(),
		authHandler: authHandler,
		roomHandler: roomHandler,
		chatHandler: chatHandler,
	}
}

type Server struct {
	httpServer *http.Server
	route      *gin.Engine

	authHandler *handlers.AuthHandler
	roomHandler *handlers.RoomHandler
	chatHandler *handlers.ChatHandler
}

func (s *Server) Run(port string) error {
	s.setupRouter()
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		MaxHeaderBytes: 1 << 20,
		Handler:        s.route,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Я не уверен, что это должно быть здесь
func (s *Server) setupRouter() {
	// auth
	s.route.POST("/sign-up", s.authHandler.SignUp)
	s.route.POST("/sign-in", s.authHandler.SignIn)
	s.route.GET("/logout", s.authHandler.Logout)

	// room
	s.route.POST("/rooms", s.authHandler.UserIdentity, s.roomHandler.PostRoom)
	s.route.GET("/rooms/:id", s.authHandler.UserIdentity, s.roomHandler.GetRoom)
	s.route.PATCH("/rooms/info", s.authHandler.UserIdentity, s.roomHandler.PatchRoomInfo)
	s.route.POST("/rooms/:id", s.authHandler.UserIdentity, s.roomHandler.DeleteRoom)

	// ws
	s.route.GET("/ws/joinRoom/:id", s.authHandler.UserIdentity, s.chatHandler.JoinRoom)
}
