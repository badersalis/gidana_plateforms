package handlers

import (
	"time"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type RentalInput struct {
	PropertyID   uint    `json:"property_id" binding:"required"`
	StartDate    string  `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate      string  `json:"end_date"`
	MonthlyPrice float64 `json:"monthly_price" binding:"required"`
}

func GetMyRentals(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var rentals []models.Rental
	database.DB.Where("tenant_id = ?", userID).
		Preload("Property.Images").
		Find(&rentals)
	utils.OK(c, rentals)
}

func CreateRental(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input RentalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var prop models.Property
	if err := database.DB.First(&prop, input.PropertyID).Error; err != nil {
		utils.NotFound(c, "Property not found")
		return
	}
	if !prop.IsAvailable {
		utils.BadRequest(c, "Property is not available")
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		utils.BadRequest(c, "Invalid start date format (YYYY-MM-DD)")
		return
	}

	rental := models.Rental{
		PropertyID:   input.PropertyID,
		TenantID:     userID,
		StartDate:    startDate,
		MonthlyPrice: input.MonthlyPrice,
		Status:       "pending",
	}

	if input.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", input.EndDate); err == nil {
			rental.EndDate = &endDate
		}
	}

	database.DB.Create(&rental)
	database.DB.Model(&prop).Update("is_available", false)

	database.DB.Preload("Property.Images").First(&rental, rental.ID)
	utils.Created(c, rental)
}

func UpdateRentalStatus(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	rentalID := c.Param("id")

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	validStatuses := map[string]bool{"pending": true, "occupied": true, "available": true, "completed": true}
	if !validStatuses[input.Status] {
		utils.BadRequest(c, "Invalid status")
		return
	}

	var rental models.Rental
	if err := database.DB.First(&rental, rentalID).Error; err != nil {
		utils.NotFound(c, "Rental not found")
		return
	}

	var prop models.Property
	database.DB.First(&prop, rental.PropertyID)
	if prop.OwnerID != userID && rental.TenantID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	database.DB.Model(&rental).Update("status", input.Status)

	if input.Status == "completed" || input.Status == "available" {
		database.DB.Model(&prop).Update("is_available", true)
	}

	utils.OK(c, gin.H{"message": "Status updated", "status": input.Status})
}
