package dashboardController

import (
	"github.com/gin-gonic/gin"
	"go-otp-auth-service/models"
	"net/http"
)

// Dashboard godoc
// @Summary Show user dashboard
// @Description Returns user info for the dashboard
// @Tags Dashboard
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /dashboard [get]
func Dashboard(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(models.User)
	c.JSON(http.StatusOK, user)
	//c.HTML(http.StatusOK, "dashboard.html", gin.H{
	//	"user": user,
	//})
}
