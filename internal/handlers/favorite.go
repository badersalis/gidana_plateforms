package handlers

import (
	"strconv"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetFavorites(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize := 10
	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&total)

	var favorites []models.Favorite
	database.DB.Where("user_id = ?", userID).
		Preload("Property.Images").
		Preload("Property.Reviews").
		Offset(offset).Limit(pageSize).
		Find(&favorites)

	props := make([]models.Property, 0, len(favorites))
	for _, f := range favorites {
		f.Property.ComputeRating()
		f.Property.IsFavorited = true
		props = append(props, f.Property)
	}

	utils.Paginated(c, props, total, page, pageSize)
}

func ToggleFavorite(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	propIDStr := c.Param("id")
	propID, err := strconv.ParseUint(propIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid property ID")
		return
	}

	var prop models.Property
	if err := database.DB.First(&prop, propID).Error; err != nil {
		utils.NotFound(c, "Property not found")
		return
	}

	var fav models.Favorite
	result := database.DB.Where("user_id = ? AND property_id = ?", userID, propID).First(&fav)

	if result.Error == nil {
		database.DB.Delete(&fav)
		utils.OK(c, gin.H{"favorited": false, "message": "Removed from favorites"})
	} else {
		newFav := models.Favorite{UserID: userID, PropertyID: uint(propID)}
		database.DB.Create(&newFav)
		utils.OK(c, gin.H{"favorited": true, "message": "Added to favorites"})
	}
}
