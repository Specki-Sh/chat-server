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
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
