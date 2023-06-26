package route

import (
	"chat-server/internal/handlers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	route       *gin.Engine
	authHandler *handlers.AuthHandler
	roomHandler *handlers.RoomHandler
	chatHandler *handlers.ChatHandler
}

func NewRouter(authHandler *handlers.AuthHandler, roomHandler *handlers.RoomHandler, chatHandler *handlers.ChatHandler) *Router {
	return &Router{
		route:       gin.New(),
		authHandler: authHandler,
		roomHandler: roomHandler,
		chatHandler: chatHandler,
	}
}

func (r *Router) SetupRouter() *gin.Engine {
	// auth
	auth := r.route.Group("/auth")
	auth.POST("/sign-up", r.authHandler.SignUp)
	auth.POST("/sign-in", r.authHandler.SignIn)
	auth.GET("/logout", r.authHandler.Logout)

	// room
	room := r.route.Group("/rooms")
	room.POST("/",
		r.authHandler.UserIdentity,
		r.roomHandler.PostRoom,
	)
	room.GET("/:id/info",
		r.roomHandler.RoomExistsMiddlewareByParam("id"),
		r.roomHandler.GetRoomInfo,
	)
	room.PATCH("/:id/info",
		r.authHandler.UserIdentity,
		r.roomHandler.RoomExistsMiddlewareByParam("id"),
		r.roomHandler.RoomPermissionsMiddleware,
		r.roomHandler.PatchRoomInfo,
	)
	room.DELETE("/:id",
		r.authHandler.UserIdentity,
		r.roomHandler.RoomExistsMiddlewareByParam("id"),
		r.roomHandler.RoomPermissionsMiddleware,
		r.roomHandler.DeleteRoom,
	)
	room.POST("/:id/members/:userID",
		r.authHandler.UserIdentity,
		r.roomHandler.RoomExistsMiddlewareByParam("id"),
		r.roomHandler.RoomPermissionsMiddleware,
		r.authHandler.UserExistMiddlewareByParam("userID"),
		r.roomHandler.AddMemberToRoomHandler,
	)
	room.DELETE("/:id/messages",
		r.authHandler.UserIdentity,
		r.roomHandler.RoomExistsMiddlewareByParam("id"),
		r.roomHandler.RoomPermissionsMiddleware,
		r.chatHandler.DeleteAllMessageFromRoom,
	)

	// chat
	r.route.GET("/chat/joinRoom/:id",
		//r.authHandler.UserIdentity,
		r.roomHandler.RoomExistsMiddlewareByParam("id"),
		//r.roomHandler.RoomAccessMiddleware,
		r.chatHandler.JoinRoom,
	)

	// message
	messages := r.route.Group("/messages")
	messages.GET("/paginate/rooms/:roomID",
		r.authHandler.UserIdentity,
		r.roomHandler.RoomExistsMiddlewareByParam("roomID"),
		r.roomHandler.RoomAccessMiddleware,
		r.chatHandler.GetMessagesPaginate,
	)
	messages.PATCH("/:id",
		r.authHandler.UserIdentity,
		r.chatHandler.MessagePermissionMiddlewareByParam("id"),
		r.chatHandler.BroadcastMessageUpdateMiddleware,
		r.chatHandler.EditMessage,
	)
	messages.DELETE("/:id",
		r.authHandler.UserIdentity,
		r.chatHandler.MessagePermissionMiddlewareByParam("id"),
		r.chatHandler.BroadcastMessageUpdateMiddleware,
		r.chatHandler.DeleteMessage,
	)

	return r.route
}
