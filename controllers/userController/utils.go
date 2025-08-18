package userController

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-otp-auth-service/initializers"
	"go-otp-auth-service/models"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"os"
	"strconv"
	"time"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashed), err
}

func UsersShow(c *gin.Context) {
	var users []models.User
	initializers.DB.Find(&users)
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "Show all users",
		"users":  users,
	})
}

func UserShow(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	initializers.DB.First(&user, id)
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "Show specific user",
		"user":   user,
	})
}

func UserUpdate(c *gin.Context) {
	var userBase struct {
		Phone string
	}
	err := c.Bind(&userBase)
	if err != nil {
		return
	}

	var user models.User
	id := c.Param("id")
	initializers.DB.First(&user, id)
	initializers.DB.Model(&user).Updates(models.User{Phone: userBase.Phone})
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "user updated successfully",
		"user":   user,
	})
}

func UserDelete(c *gin.Context) {
	id := c.Param("id")
	initializers.DB.Delete(&models.User{}, id)
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "user deleted successfully",
	})
}

func GenerateOTP(phone string) (string, error) {
	otpExpMin, err := strconv.Atoi(os.Getenv("OTP_EXP_MIN"))
	if err != nil {
		otpExpMin = 1
	}
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	otp := fmt.Sprintf("%06d", n.Int64())
	err = initializers.RDB.Set(initializers.Ctx, "otp:"+phone, otp, time.Duration(otpExpMin)*time.Minute).Err()
	if err != nil {
		return "", err
	}
	return otp, nil
}

func VerifyOTP(phone string, otp string) bool {
	storedOtp, err := initializers.RDB.Get(initializers.Ctx, "otp:"+phone).Result()
	if err != nil {
		return false
	}
	if storedOtp != otp {
		return false
	}
	initializers.RDB.Del(initializers.Ctx, "otp:"+phone)
	return true
}
