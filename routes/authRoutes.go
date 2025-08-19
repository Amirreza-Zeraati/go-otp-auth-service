package routes

import (
	"github.com/gin-gonic/gin"
	"go-otp-auth-service/controllers/userController"
)

func AuthRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "auth.html", nil)
	})
	r.POST("/login", controllers.Login)
	r.POST("/request-otp", controllers.RequestOTP)
	r.GET("/logout", controllers.Logout)
}
