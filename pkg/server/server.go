package server

import (
	"chat-server/internal/handlers"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func NewServer(authHandler *handlers.AuthHandler) *Server {
	return &Server{
		httpServer:  &http.Server{},
		route:       gin.New(),
		authHandler: authHandler}
}

type Server struct {
	httpServer *http.Server
	route      *gin.Engine

	authHandler *handlers.AuthHandler
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
	s.route.POST("/sign-up", s.authHandler.SignUp)
	s.route.POST("/sign-in", s.authHandler.SignIn)
	s.route.GET("/logout", s.authHandler.Logout)
}
