package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}

func getUsername(c *gin.Context) (string, error) {
	name, ok := c.Get(usernameCtx)
	if !ok {
		return "", errors.New("username not found")
	}

	nameString, ok := name.(string)
	if !ok {
		return "", errors.New("username is of invalid type")
	}

	return nameString, nil
}
