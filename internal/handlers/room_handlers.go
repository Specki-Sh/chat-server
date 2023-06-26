package handlers

import (
	"net/http"
	"strconv"

	"chat-server/internal/domain/use_case"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomUseCase use_case.RoomUseCase
}

func NewRoomHandler(roomUseCase use_case.RoomUseCase) *RoomHandler {
	return &RoomHandler{
		roomUseCase: roomUseCase,
	}
}

func (r *RoomHandler) PostRoom(c *gin.Context) {
	var req use_case.CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.OwnerID = id

	res, err := r.roomUseCase.CreateRoom(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *RoomHandler) GetRoomInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := r.roomUseCase.GetRoomInfoByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (r *RoomHandler) PatchRoomInfo(c *gin.Context) {
	var req use_case.EditRoomReq
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = roomID
	res, err := r.roomUseCase.EditRoomInfo(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *RoomHandler) DeleteRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.roomUseCase.RemoveRoomByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (r *RoomHandler) AddMemberToRoomHandler(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	member, err := r.roomUseCase.AddMemberToRoom(roomID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (r *RoomHandler) RoomExistsMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		exists, err := r.roomUseCase.RoomExists(roomID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	}
}

func (r *RoomHandler) RoomExistsMiddlewareByJSON(jsonKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonInput map[string]interface{}
		if err := c.ShouldBindJSON(&jsonInput); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		roomID, ok := jsonInput[jsonKey].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
			return
		}
		exists, err := r.roomUseCase.RoomExists(int(roomID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	}
}

func (r *RoomHandler) RoomPermissionsMiddleware(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isOwner, err := r.roomUseCase.IsRoomOwner(roomID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwner {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.Next()
}

func (r *RoomHandler) RoomAccessMiddleware(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasAccess, err := r.roomUseCase.HasRoomAccess(roomID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !hasAccess {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.Next()
}
