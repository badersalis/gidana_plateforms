package handlers

import (
	"fmt"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	appws "github.com/badersalis/gidana_backend/internal/ws"
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

// notifyMatchingAlerts finds all active alerts that match prop and notifies
// their owners via WebSocket (if online) or Expo push (if offline).
// It is called asynchronously after a new property is created.
func notifyMatchingAlerts(prop models.Property) {
	var alerts []models.Alert
	database.DB.Where(
		`is_active = true
		 AND user_id != ?
		 AND (neighborhood = '' OR neighborhood = ?)
		 AND (property_type = '' OR property_type = ?)
		 AND (transaction_type = '' OR transaction_type = ?)
		 AND (min_rooms = 0 OR min_rooms <= ?)
		 AND (max_price = 0 OR max_price >= ?)`,
		prop.OwnerID,
		prop.Neighborhood,
		prop.PropertyType,
		prop.TransactionType,
		prop.Rooms,
		prop.Price,
	).Find(&alerts)

	if len(alerts) == 0 {
		return
	}

	title := "Nouvelle propriété disponible"
	body := fmt.Sprintf("%s à %s — %s (%.0f %s)",
		prop.PropertyType, prop.Neighborhood, prop.TransactionType, prop.Price, prop.Currency)
	data := map[string]any{"property_id": prop.ID}

	event := appws.Event{Type: "property_alert", Data: map[string]any{
		"property_id":      prop.ID,
		"title":            prop.Title,
		"neighborhood":     prop.Neighborhood,
		"property_type":    prop.PropertyType,
		"transaction_type": prop.TransactionType,
		"price":            prop.Price,
		"currency":         prop.Currency,
	}}

	for _, alert := range alerts {
		appws.H.Emit(alert.UserID, event)

		if !appws.H.IsOnline(alert.UserID) {
			var user models.User
			if err := database.DB.Select("expo_push_token").First(&user, alert.UserID).Error; err == nil {
				utils.SendExpoPush(user.ExpoPushToken, title, body, data)
			}
		}
	}
}
