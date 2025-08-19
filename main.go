package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-otp-auth-service/controllers/dashboardController"
	_ "go-otp-auth-service/docs"
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run()
	if err != nil {
		return
	}
}
