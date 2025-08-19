package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-otp-auth-service/dto"
	"go-otp-auth-service/initializers"
	"go-otp-auth-service/models"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Login godoc
// @Summary Login with phone and OTP
// @Description Verify OTP and login the user
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param body body dto.LoginRequest true "Login info"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var body dto.LoginRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if !VerifyOTP(body.Phone, body.OTP) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired OTP",
		})
		return
	}
	var user models.User
	initializers.DB.First(&user, "phone = ?", body.Phone)
	if user.ID == 0 {
		newUser := models.User{Phone: body.Phone}
		result := initializers.DB.Create(&newUser)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create user",
			})
			return
		}
		user = newUser
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	secret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 3600, "", "", false, true)
	c.JSON(http.StatusOK, dto.LoginResponse{Token: tokenString})
	//c.Redirect(http.StatusFound, "/dashboard")
}

// Logout godoc
// @Summary Logout user
// @Description Clear token cookie and logout
// @Tags Auth
// @Produce  json
// @Success 200 {object} dto.LogoutResponse
// @Router /logout [get]
func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, dto.LogoutResponse{Message: "Logged out successfully"})
	//c.Redirect(http.StatusFound, "/")
}

// RequestOTP godoc
// @Summary Request OTP
// @Description Generate and send OTP for login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param body body dto.RequestOTPRequest true "Phone number"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /request-otp [post]
func RequestOTP(c *gin.Context) {
	rateLimit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	periodTime, _ := strconv.Atoi(os.Getenv("PERIOD_TIME"))

	var body dto.RequestOTPRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	ok, msg, err := CheckOTPRequest(body.Phone, rateLimit, time.Duration(periodTime)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error - (first run redis)",
		})
		return
	}
	if !ok {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": msg,
		})
		return
	}
	otp, err := GenerateOTP(body.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println("otp : ", otp)
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully",
	})
}
