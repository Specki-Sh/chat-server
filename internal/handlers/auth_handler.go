package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
	usernameCtx         = "username"
)

type AuthHandler struct {
	userUseCase  use_case.UserUseCase
	authUseCase  use_case.AuthUseCase
	tokenUseCase use_case.TokenUseCase

	logger *logrus.Logger
}

func NewAuthHandler(
	uus use_case.UserUseCase,
	aus use_case.AuthUseCase,
	tus use_case.TokenUseCase,
	logger *logrus.Logger,
) *AuthHandler {
	return &AuthHandler{
		userUseCase:  uus,
		authUseCase:  aus,
		tokenUseCase: tus,
		logger:       logger,
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
	var req entity.SignInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Errorf("error binding JSON: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		a.logger.Errorf("Error sign-in user data: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.authUseCase.Authenticate(&req)
	if err != nil {
		a.logger.Errorf("error generating token: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	a.logger.Infof("user with email %s signed in", req.Email)
	c.JSON(http.StatusOK, res)
}

func (a *AuthHandler) Logout(c *gin.Context) {
	var req entity.LogoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Errorf("error binding JSON: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.RefreshToken == "" {
		a.logger.Error("error empty refresh token")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error empty refresh token"})
		return
	}

	if err := a.authUseCase.Logout(c, req.RefreshToken); err != nil {
		a.logger.Errorf("error invalidate token: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "empty auth header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "invalid auth header"})
		return
	}

	if len(headerParts[1]) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "token is empty"})
		return
	}

	userID, username, err := a.tokenUseCase.ParseAccessToken(headerParts[1])
	if err != nil {
		a.logger.Errorf("error parsing token: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Set(userCtx, userID)
	c.Set(usernameCtx, username)
	a.logger.Infof("user identity set: %d %s", userID, username)
}

func (a *AuthHandler) UserIdentityByQueryParam(query string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query(query)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"reason": "empty access token"})
			return
		}

		userID, username, err := a.tokenUseCase.ParseAccessToken(token)
		if err != nil {
			a.logger.Errorf("error parsing token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set(userCtx, userID)
		c.Set(usernameCtx, username)
		a.logger.Infof("user identity set: %d %s", userID, username)
	}
}

func (a *AuthHandler) UserPermissionMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInt, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			a.logger.Errorf("Error converting room ID to int: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}
		userID := entity.ID(userIDInt)

		userTokenID, err := getUserID(c)
		if err != nil {
			a.logger.Errorf("\"Error getting user ID: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if userID != userTokenID {
			a.logger.Infof(
				"User with ID: %v has not permission of user with ID: %v",
				userTokenID,
				userID,
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.Next()
	}
}

func (a *AuthHandler) UserExistMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInt, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			a.logger.Errorf("error converting userID to int: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID := entity.ID(userIDInt)
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
