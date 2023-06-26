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

func (h *RoomHandler) PostRoom(c *gin.Context) {
	var req use_case.CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.OwnerId = id

	res, err := h.roomUseCase.CreateRoom(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.roomUseCase.GetRoomByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) PatchRoomInfo(c *gin.Context) {
	var req use_case.EditRoomReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.roomUseCase.EditRoomInfo(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.roomUseCase.RemoveRoomByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *RoomHandler) AddMemberToRoomHandler(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	member, err := h.roomUseCase.AddMemberToRoom(roomID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (h *RoomHandler) RoomExistsMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		exists, err := h.roomUseCase.RoomExists(roomID)
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

func (h *RoomHandler) RoomExistsMiddlewareByJSON(jsonKey string) gin.HandlerFunc {
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
		exists, err := h.roomUseCase.RoomExists(int(roomID))
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

func (h *RoomHandler) RoomPermissionsMiddleware(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isOwner, err := h.roomUseCase.IsRoomOwner(roomID, userID)
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

func (h *RoomHandler) RoomAccessMiddleware(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasAccess, err := h.roomUseCase.HasRoomAccess(roomID, userID)
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
