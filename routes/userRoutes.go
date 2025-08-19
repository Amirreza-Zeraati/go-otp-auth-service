package routes

import (
	"github.com/gin-gonic/gin"
	"go-otp-auth-service/controllers/userController"
)

func UserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("/:id", controllers.GetUser)
		users.GET("", controllers.ListUsers)
	}
}
