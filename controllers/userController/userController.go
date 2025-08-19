package controllers

import (
	"github.com/gin-gonic/gin"
	"go-otp-auth-service/initializers"
	"go-otp-auth-service/models"
	"net/http"
	"strconv"
)

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.HTML(http.StatusOK, "user_detail.html", gin.H{
		"User": user,
	})
}

func ListUsers(c *gin.Context) {
	var users []models.User
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit := 5
	offset := (page - 1) * limit

	search := c.Query("search")

	query := initializers.DB.Model(&models.User{})
	if search != "" {
		query = query.Where("phone LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	query.Limit(limit).Offset(offset).Order("created_at desc").Find(&users)

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.HTML(http.StatusOK, "users.html", gin.H{
		"Users":      users,
		"Page":       page,
		"PrevPage":   page - 1,
		"NextPage":   page + 1,
		"HasPrev":    page > 1,
		"HasNext":    page < totalPages,
		"TotalPages": totalPages,
		"Search":     search,
	})

}
