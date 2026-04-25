package handlers

import (
	"strings"
	"time"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetSearchSuggestions(c *gin.Context) {
	q := c.Query("q")
	if len(q) < 2 {
		utils.OK(c, []string{})
		return
	}

	like := "%" + strings.ToLower(q) + "%"
	var neighborhoods []struct{ Neighborhood string }
	database.DB.Model(&models.Property{}).
		Select("DISTINCT neighborhood").
		Where("LOWER(neighborhood) LIKE ?", like).
		Limit(10).
		Scan(&neighborhoods)

	suggestions := make([]string, 0, len(neighborhoods))
	for _, n := range neighborhoods {
		suggestions = append(suggestions, n.Neighborhood)
	}

	utils.OK(c, suggestions)
}

func SaveSearchHistory(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.OK(c, gin.H{"saved": false})
		return
	}

	var input struct {
		SearchTerm string `json:"search_term" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	oneHourAgo := time.Now().Add(-time.Hour)
	var existing models.SearchHistory
	result := database.DB.Where("user_id = ? AND search_term = ? AND created_at > ?",
		userID, input.SearchTerm, oneHourAgo).First(&existing)

	if result.Error != nil {
		sh := models.SearchHistory{UserID: userID, SearchTerm: input.SearchTerm}
		database.DB.Create(&sh)
	}

	utils.OK(c, gin.H{"saved": true})
}

func GetSearchHistory(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var history []models.SearchHistory
	database.DB.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(10).
		Find(&history)
	utils.OK(c, history)
}

func DeleteSearchHistory(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	database.DB.Where("user_id = ?", userID).Delete(&models.SearchHistory{})
	utils.OK(c, gin.H{"message": "Search history cleared"})
}
