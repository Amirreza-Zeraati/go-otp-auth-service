package dashboardController

import (
	"github.com/gin-gonic/gin"
	"go-otp-auth-service/models"
	"net/http"
)

func Dashboard(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(models.User)
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"user": user,
	})
}
