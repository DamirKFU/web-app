package controller

import (
	"app/internal/core"
	"app/internal/models"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	server *core.Server
}

func NewUsersController(server *core.Server) *UsersController {
	return &UsersController{server: server}
}

func (ctrl *UsersController) Create(ctx *gin.Context) {
	var lastUser models.User
	result := ctrl.server.DB.Order("CAST(SUBSTR(username, 9) AS INTEGER) DESC").First(&lastUser)

	newUsername := "username1"
	if result.Error == nil {
		re := regexp.MustCompile(`^username(\d+)$`)
		matches := re.FindStringSubmatch(lastUser.Username)
		if len(matches) == 2 {
			lastNum, _ := strconv.Atoi(matches[1])
			newUsername = "username" + strconv.Itoa(lastNum+1)
		}
	}

	newUser := models.User{Username: newUsername, Password: "123"}
	if err := ctrl.server.DB.Create(&newUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать пользователя"})
		return
	}

	var usernames []string
	if err := ctrl.server.DB.Model(&models.User{}).Order("CAST(SUBSTR(username, 9) AS INTEGER)").Pluck("username", &usernames).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить список пользователей"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"users":  usernames,
	})
}
