package main

import (
	"github.com/gin-gonic/gin"
	"go-otp-auth-service/controllers/dashboardController"
	"go-otp-auth-service/controllers/userController"
	"go-otp-auth-service/initializers"
	middleware "go-otp-auth-service/middleware/auth"
)

func init() {
	initializers.LoadEnvFile()
	initializers.ConnectToDB()
	initializers.Migrate()
	initializers.ConnectRedis()
}

func main() {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "auth.html", nil)
	})
	r.POST("/request-otp", userController.RequestOTP)
	r.POST("/login", userController.Login)
	r.GET("/logout", userController.Logout)

	r.GET("/dashboard", middleware.RequireAuth, dashboardController.Dashboard)

	err := r.Run()
	if err != nil {
		return
	}
}
