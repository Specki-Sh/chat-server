package handlers

import (
	"chat-server/internal/domain/entity"
	"errors"
	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) (entity.ID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return entity.ID(idInt), nil
}

func getUsername(c *gin.Context) (entity.NonEmptyString, error) {
	name, ok := c.Get(usernameCtx)
	if !ok {
		return "", errors.New("username not found")
	}

	nameString, ok := name.(string)
	if !ok {
		return "", errors.New("username is of invalid type")
	}

	return entity.NonEmptyString(nameString), nil
}
