package main

import (
	"github.com/gin-gonic/gin"
	"go-otp-auth-service/controllers/dashboardController"
	"go-otp-auth-service/initializers"
	middleware "go-otp-auth-service/middleware/auth"
	"go-otp-auth-service/routes"
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

	routes.AuthRoutes(r)
	routes.UserRoutes(r)

	r.GET("/dashboard", middleware.RequireAuth, dashboardController.Dashboard)

	err := r.Run()
	if err != nil {
		return
	}
}
