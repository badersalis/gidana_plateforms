package handlers

import (
	"time"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type UpdateProfileInput struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Locale      string `json:"locale"`
	Timezone    string `json:"timezone"`
}

func UpdateProfile(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if input.FirstName != "" {
		updates["first_name"] = input.FirstName
	}
	if input.LastName != "" {
		updates["last_name"] = input.LastName
	}
	if input.Gender != "" {
		updates["gender"] = input.Gender
	}
	if input.Locale != "" {
		updates["locale"] = input.Locale
	}
	if input.Timezone != "" {
		updates["timezone"] = input.Timezone
	}
	if input.DateOfBirth != "" {
		t, err := time.Parse("2006-01-02", input.DateOfBirth)
		if err == nil {
			updates["date_of_birth"] = t
		}
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update profile")
		return
	}

	var user models.User
	database.DB.First(&user, userID)
	utils.OK(c, user)
}

func UploadProfilePicture(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	file, err := c.FormFile("picture")
	if err != nil {
		utils.BadRequest(c, "No file provided")
		return
	}

	url, err := saveFile(c, file)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("profile_picture", url)

	var user models.User
	database.DB.First(&user, userID)
	utils.OK(c, user)
}

func UpdatePushToken(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		Token string `json:"expo_push_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("expo_push_token", input.Token)
	utils.OK(c, gin.H{"message": "Push token updated"})
}

func ChangePassword(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var user models.User
	database.DB.First(&user, userID)

	if !utils.CheckPassword(input.CurrentPassword, user.PasswordHash) {
		utils.BadRequest(c, "Current password is incorrect")
		return
	}

	hash, _ := utils.HashPassword(input.NewPassword)
	database.DB.Model(&user).Update("password_hash", hash)
	utils.OK(c, gin.H{"message": "Password changed successfully"})
}
