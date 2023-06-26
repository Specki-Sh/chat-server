package handlers

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	logger *logrus.Logger
}

func NewAuthHandler(uus use_case.UserUseCase, aus use_case.AuthUseCase, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		userUseCase: uus,
		authUseCase: aus,
		logger:      logger,
	}
}

func (a *AuthHandler) SignUp(c *gin.Context) {
	var u entity.CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		a.logger.Errorf("error binding JSON: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.userUseCase.CreateUser(&u)
	if err != nil {
		a.logger.Errorf("error creating user: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	a.logger.Infof("user created: %v", res)
	c.JSON(http.StatusOK, res)
}

func (a *AuthHandler) SignIn(c *gin.Context) {
	var user entity.SignInReq
	if err := c.ShouldBindJSON(&user); err != nil {
		a.logger.Errorf("error binding JSON: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := a.authUseCase.GenerateToken(&user)
	if err != nil {
		a.logger.Errorf("error generating token: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", u.AccessToken, 60*60*24, "/", "localhost", false, true)
	a.logger.Infof("user signed in: %v", u)
	c.JSON(http.StatusOK, u)
}

func (a *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	a.logger.Info("user logged out")
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (a *AuthHandler) UserIdentity(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		a.logger.Errorf("error getting jwt cookie: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "no jwt cookie"})
		return
	}

	if len(cookie) == 0 {
		a.logger.Error("jwt cookie is empty")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "token is empty"})
		return
	}

	userID, username, err := a.authUseCase.ParseToken(cookie)
	if err != nil {
		a.logger.Errorf("error parsing token: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": err.Error()})
		return
	}

	c.Set(userCtx, userID)
	c.Set(usernameCtx, username)
	a.logger.Infof("user identity set: %d %s", userID, username)
}

func (a *AuthHandler) UserExistMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			a.logger.Errorf("error converting userID to int: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		exists, err := a.userUseCase.UserExists(userID)
		if err != nil {
			a.logger.Errorf("error checking if user exists: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			a.logger.Infof("user not found: %d", userID)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		a.logger.Infof("user exists: %d", userID)
		c.Next()
	}
}
