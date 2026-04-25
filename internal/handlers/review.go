package handlers

import (
	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type ReviewInput struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

func GetPropertyReviews(c *gin.Context) {
	propID := c.Param("id")
	var reviews []models.Review
	database.DB.Where("property_id = ?", propID).Preload("User").Find(&reviews)
	utils.OK(c, reviews)
}

func CreateReview(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	propID := c.Param("id")

	var prop models.Property
	if err := database.DB.First(&prop, propID).Error; err != nil {
		utils.NotFound(c, "Property not found")
		return
	}

	var existing models.Review
	if err := database.DB.Where("user_id = ? AND property_id = ?", userID, propID).First(&existing).Error; err == nil {
		utils.BadRequest(c, "You have already reviewed this property")
		return
	}

	var input ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	review := models.Review{
		PropertyID: prop.ID,
		UserID:     userID,
		Rating:     input.Rating,
		Comment:    input.Comment,
	}

	database.DB.Create(&review)
	database.DB.Preload("User").First(&review, review.ID)
	utils.Created(c, review)
}

func DeleteReview(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	reviewID := c.Param("id")

	var review models.Review
	if err := database.DB.First(&review, reviewID).Error; err != nil {
		utils.NotFound(c, "Review not found")
		return
	}
	if review.UserID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	database.DB.Delete(&review)
	utils.OK(c, gin.H{"message": "Review deleted"})
}
