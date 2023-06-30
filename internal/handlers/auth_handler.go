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
	var user entity.CreateUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		a.logger.Errorf("error binding JSON: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := user.Validate(); err != nil {
		a.logger.Errorf("Error create user data: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.userUseCase.CreateUser(&user)
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
	if err := user.Validate(); err != nil {
		a.logger.Errorf("Error sign-in user data: %v", err)
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

func (a *AuthHandler) PatchUserProfile(c *gin.Context) {
	var req entity.EditProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Errorf("error binding JSON: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		a.logger.Errorf("Error user profile data: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.userUseCase.EditUserProfile(&req)
	if err != nil {
		a.logger.Errorf("error editing user profile: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	a.logger.Infof("user %v changed profile %v", req.ID, req)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
	c.JSON(http.StatusOK, user)
}

func (a *AuthHandler) UserIdentity(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		a.logger.Errorf("error getting jwt cookie: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no jwt cookie"})
		return
	}

	if len(cookie) == 0 {
		a.logger.Error("jwt cookie is empty")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
		return
	}

	userID, username, err := a.authUseCase.ParseToken(cookie)
	if err != nil {
		a.logger.Errorf("error parsing token: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Set(userCtx, userID)
	c.Set(usernameCtx, username)
	a.logger.Infof("user identity set: %d %s", userID, username)
}

func (a *AuthHandler) UserPermissionMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			a.logger.Errorf("Error converting room ID to int: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		userTokenID, err := getUserID(c)
		if err != nil {
			a.logger.Errorf("\"Error getting user ID: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if userID != userTokenID {
			a.logger.Infof("User with ID: %v has not permission of user with ID: %v", userTokenID, userID)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.Next()
	}
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
