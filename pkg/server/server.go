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
	auth := s.route.Group("/auth")
	auth.POST("/sign-up", s.authHandler.SignUp)
	auth.POST("/sign-in", s.authHandler.SignIn)
	auth.GET("/logout", s.authHandler.Logout)

	// room
	room := s.route.Group("/rooms")
	room.POST("/", s.roomHandler.PostRoom)
	room.GET("/:id", s.roomHandler.GetRoom)
	room.PATCH("/info", s.roomHandler.PatchRoomInfo)
	room.POST("/:id", s.roomHandler.DeleteRoom)

	// ws
	ws := s.route.Group("/ws")
	ws.GET("/ws/joinRoom/:id", s.chatHandler.JoinRoom)
}
