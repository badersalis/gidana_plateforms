package handlers

import (
	"time"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func RequestDeleteAccount(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	// Prevent duplicate requests
	var existing models.DeletedAccount
	if err := database.DB.Where("user_id = ?", userID).First(&existing).Error; err == nil {
		utils.BadRequest(c, "Account deletion already requested")
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.InternalError(c, "Failed to process deletion request")
		return
	}

	snapshot := models.DeletedAccount{
		UserID:         user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: user.ProfilePicture,
		Gender:         user.Gender,
		DateOfBirth:    user.DateOfBirth,
		MemberSince:    user.MemberSince,
		Locale:         user.Locale,
		Timezone:       user.Timezone,
		RequestedAt:    time.Now(),
		Status:         "pending",
	}

	if err := database.DB.Create(&snapshot).Error; err != nil {
		utils.InternalError(c, "Failed to process deletion request")
		return
	}

	// Deactivate immediately so existing tokens stop working, then soft-delete
	database.DB.Model(&user).Update("active", false)
	database.DB.Delete(&user)

	utils.SendExpoPush(
		user.ExpoPushToken,
		"Deletion Request Received",
		"Your account deletion request has been received. Our compliance team will review it and notify you once the process is complete.",
		gin.H{"type": "account_deletion_requested"},
	)

	utils.OK(c, gin.H{"message": "Account deletion request submitted. You will be notified once reviewed."})
}
