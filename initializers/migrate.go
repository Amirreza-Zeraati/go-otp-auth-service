package initializers

import (
	"go-otp-auth-service/models"
	"log"
)

func Migrate() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}
	log.Println("Database migrated successfully")
}
