package handlers

import (
	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type AlertInput struct {
	Neighborhood    string  `json:"neighborhood"`
	PropertyType    string  `json:"property_type"`
	MinRooms        int     `json:"min_rooms"`
	MaxPrice        float64 `json:"max_price"`
	TransactionType string  `json:"transaction_type"`
}

func GetAlerts(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var alerts []models.Alert
	database.DB.Where("user_id = ?", userID).Find(&alerts)
	utils.OK(c, alerts)
}

func CreateAlert(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input AlertInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	alert := models.Alert{
		UserID:          userID,
		Neighborhood:    input.Neighborhood,
		PropertyType:    input.PropertyType,
		MinRooms:        input.MinRooms,
		MaxPrice:        input.MaxPrice,
		TransactionType: input.TransactionType,
		IsActive:        true,
	}

	database.DB.Create(&alert)
	utils.Created(c, alert)
}

func UpdateAlert(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	alertID := c.Param("id")

	var alert models.Alert
	if err := database.DB.Where("id = ? AND user_id = ?", alertID, userID).First(&alert).Error; err != nil {
		utils.NotFound(c, "Alert not found")
		return
	}

	var input struct {
		AlertInput
		IsActive *bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{
		"neighborhood":     input.Neighborhood,
		"property_type":    input.PropertyType,
		"min_rooms":        input.MinRooms,
		"max_price":        input.MaxPrice,
		"transaction_type": input.TransactionType,
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	database.DB.Model(&alert).Updates(updates)
	utils.OK(c, alert)
}

func DeleteAlert(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	alertID := c.Param("id")

	var alert models.Alert
	if err := database.DB.Where("id = ? AND user_id = ?", alertID, userID).First(&alert).Error; err != nil {
		utils.NotFound(c, "Alert not found")
		return
	}

	database.DB.Delete(&alert)
	utils.OK(c, gin.H{"message": "Alert deleted"})
}
