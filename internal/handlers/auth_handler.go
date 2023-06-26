package handlers

import (
	entity2 "chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	userCtx     = "userID"
	usernameCtx = "username"
)

type AuthHandler struct {
	userUseCase use_case.UserUseCase
	authUseCase use_case.AuthUseCase
}

func NewAuthHandler(uus use_case.UserUseCase, aus use_case.AuthUseCase) *AuthHandler {
	return &AuthHandler{userUseCase: uus, authUseCase: aus}
}

func (a *AuthHandler) SignUp(c *gin.Context) {
	var u entity2.CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.userUseCase.CreateUser(&u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *AuthHandler) SignIn(c *gin.Context) {
	var user entity2.SignInReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := a.authUseCase.GenerateToken(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", u.AccessToken, 60*60*24, "/", "localhost", false, true)
	c.JSON(http.StatusOK, u)
}

func (a *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (a *AuthHandler) UserIdentity(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "no jwt cookie"})
		return
	}

	if len(cookie) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "token is empty"})
		return
	}

	userID, username, err := a.authUseCase.ParseToken(cookie)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": err.Error()})
		return
	}

	c.Set(userCtx, userID)
	c.Set(usernameCtx, username)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (a *AuthHandler) UserExistMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		exists, err := a.userUseCase.UserExists(userID)
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
