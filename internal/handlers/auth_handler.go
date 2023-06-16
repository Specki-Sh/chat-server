package handlers

import (
	entity2 "chat-server/internal/domain/entity"
	use_case2 "chat-server/internal/domain/use_case"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userUseCase use_case2.UserUseCase
	authUseCase use_case2.AuthUseCase
}

func NewAuthHandler(uus use_case2.UserUseCase, aus use_case2.AuthUseCase) *AuthHandler {
	return &AuthHandler{userUseCase: uus, authUseCase: aus}
}

func (a *AuthHandler) SignUp(c *gin.Context) {
	var u entity2.CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.userUseCase.CreateUser(&u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *AuthHandler) SignIn(c *gin.Context) {
	var user entity2.SignInReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := a.authUseCase.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", u.AccessToken, 60*60*24, "/", "localhost", false, true)
	c.JSON(http.StatusOK, u)
}

func (a *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
