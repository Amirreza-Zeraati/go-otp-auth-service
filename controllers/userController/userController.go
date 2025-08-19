package controllers

import (
	"github.com/gin-gonic/gin"
	"go-otp-auth-service/dto"
	"go-otp-auth-service/initializers"
	"go-otp-auth-service/models"
	"net/http"
	"strconv"
)

// GetUser godoc
// GetUser @Summary Get user by ID
// @Description Returns a single user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	response := dto.UserResponse{
		ID:        user.ID,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusOK, response)

	//c.HTML(http.StatusOK, "user.html", gin.H{"User": user})
}

// ListUsers godoc
// @Summary List users
// @Description Returns paginated list of users (with optional search)
// @Tags users
// @Produce json
// @Param page query int false "Page number"
// @Param search query string false "Search by phone"
// @Success 200 {object} dto.UsersListResponse
// @Router /users [get]
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
	var userResponses []dto.UserResponse
	for _, u := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        u.ID,
			Phone:     u.Phone,
			CreatedAt: u.CreatedAt,
		})
	}
	response := dto.UsersListResponse{
		Users:      userResponses,
		Page:       page,
		PrevPage:   page - 1,
		NextPage:   page + 1,
		HasPrev:    page > 1,
		HasNext:    page < totalPages,
		TotalPages: totalPages,
		Search:     search,
	}
	c.JSON(http.StatusOK, response)

	//c.HTML(http.StatusOK, "users.html", gin.H{
	//	"Users":      users,
	//	"Page":       page,
	//	"PrevPage":   page - 1,
	//	"NextPage":   page + 1,
	//	"HasPrev":    page > 1,
	//	"HasNext":    page < totalPages,
	//	"TotalPages": totalPages,
	//	"Search":     search,
	//})
}
