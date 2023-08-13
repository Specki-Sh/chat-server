package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

type RoomHandler struct {
	roomUseCase use_case.RoomUseCase

	logger *logrus.Logger
}

func NewRoomHandler(roomUseCase use_case.RoomUseCase, logger *logrus.Logger) *RoomHandler {
	return &RoomHandler{
		roomUseCase: roomUseCase,
		logger:      logger,
	}
}

func (r *RoomHandler) PostRoom(c *gin.Context) {
	var req entity.CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		r.logger.Errorf("Error binding JSON: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := getUserID(c)
	if err != nil {
		r.logger.Errorf("Error getting user ID: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.OwnerID = id

	res, err := r.roomUseCase.CreateRoom(&req)
	if err != nil {
		r.logger.Errorf("Error creating room: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r.logger.Infof("Successfully created room with ID: %v", res.ID)
	c.JSON(http.StatusOK, res)
}

func (r *RoomHandler) GetRoomInfo(c *gin.Context) {
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.logger.Errorf("Error converting room ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := entity.ID(idInt)

	room, err := r.roomUseCase.GetRoomInfoByID(id)
	if err != nil {
		r.logger.Errorf("Error getting room info by ID: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r.logger.Infof("Successfully got room info for ID: %v", id)
	c.JSON(http.StatusOK, room)
}

func (r *RoomHandler) PatchRoomInfo(c *gin.Context) {
	var req entity.EditRoomReq
	if err := c.BindJSON(&req); err != nil {
		r.logger.Errorf("Error binding JSON: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomIDInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.logger.Errorf("Error converting room ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = entity.ID(roomIDInt)
	res, err := r.roomUseCase.EditRoomInfo(&req)
	if err != nil {
		r.logger.Errorf("Error editing room info: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r.logger.Infof("Successfully edited room info for ID: %v", roomIDInt)
	c.JSON(http.StatusOK, res)
}

func (r *RoomHandler) DeleteRoom(c *gin.Context) {
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.logger.Errorf("Error converting room ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := entity.ID(idInt)

	err = r.roomUseCase.RemoveRoomByID(id)
	if err != nil {
		r.logger.Errorf("Error removing room by ID: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r.logger.Infof("Successfully deleted room with ID: %v", id)
	c.Status(http.StatusNoContent)
}

func (r *RoomHandler) AddMemberToRoomHandler(c *gin.Context) {
	roomIDInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.logger.Errorf("Error converting room ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	userIDInt, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		r.logger.Errorf("Error converting user ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	roomID, userID := entity.ID(roomIDInt), entity.ID(userIDInt)

	member, err := r.roomUseCase.AddMemberToRoom(roomID, userID)
	if err != nil {
		r.logger.Errorf("Error adding member to room: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	r.logger.Infof("Successfully added member with ID: %v to room with ID: %v", userID, roomID)
	c.JSON(http.StatusOK, member)
}

func (r *RoomHandler) RoomExistsMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomIDInt, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			r.logger.Errorf("Error converting room ID to int: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		roomID := entity.ID(roomIDInt)
		exists, err := r.roomUseCase.RoomExists(roomID)
		if err != nil {
			r.logger.Errorf("Error checking if room exists: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			r.logger.Infof("Room with ID: %v does not exist", roomID)
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
			r.logger.Errorf("Error binding JSON: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		roomID, ok := jsonInput[jsonKey].(float64)
		if !ok {
			r.logger.Infof("Invalid room ID in JSON input")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
			return
		}
		exists, err := r.roomUseCase.RoomExists(entity.ID(roomID))
		if err != nil {
			r.logger.Errorf("Error checking if room exists: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			r.logger.Infof("Room with ID: %v does not exist", int(roomID))
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	}
}

func (r *RoomHandler) RoomPermissionsMiddleware(c *gin.Context) {
	roomIDInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.logger.Errorf("Error converting room ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}
	roomID := entity.ID(roomIDInt)

	userID, err := getUserID(c)
	if err != nil {
		r.logger.Errorf("Error getting user ID: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isOwner, err := r.roomUseCase.IsRoomOwner(roomID, userID)
	if err != nil {
		r.logger.Errorf("Error checking if user is owner of room: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isOwner {
		r.logger.Infof("User with ID: %v is not owner of room with ID: %v", userID, roomID)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	c.Next()
}

func (r *RoomHandler) RoomAccessMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomIDInt, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			r.logger.Errorf("Error converting room ID to int: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
			return
		}
		roomID := entity.ID(roomIDInt)

		userID, err := getUserID(c)
		if err != nil {
			r.logger.Errorf("Error getting user ID: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hasAccess, err := r.roomUseCase.HasRoomAccess(roomID, userID)
		if err != nil {
			r.logger.Errorf("Error checking if user has access to room: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasAccess {
			r.logger.Infof(
				"User with ID: %v does not have access to room with ID: %v",
				userID,
				roomID,
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.Next()
	}
}
