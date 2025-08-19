package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-otp-auth-service/initializers"
	"go-otp-auth-service/models"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Login(c *gin.Context) {
	var body struct {
		Phone string
		OTP   string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if !VerifyOTP(body.Phone, body.OTP) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
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
	c.Redirect(http.StatusFound, "/dashboard")
}

func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func RequestOTP(c *gin.Context) {
	rateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	periodTime, err := strconv.Atoi(os.Getenv("PERIOD_TIME"))
	var body struct {
		Phone string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	ok, msg, err := CheckOTPRequest(body.Phone, rateLimit, time.Duration(periodTime)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error - (first run redis)"})
		return
	}
	if !ok {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": msg})
		return
	}
	otp, err := GenerateOTP(body.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("otp : ", otp)
}
