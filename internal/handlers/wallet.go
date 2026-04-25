package handlers

import (
	"strings"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type WalletInput struct {
	Provider       string `json:"provider" binding:"required"`
	Nature         string `json:"nature"`
	PhoneNumber    string `json:"phone_number"`
	Email          string `json:"email"`
	CardNumber     string `json:"card_number"`
	CVV            string `json:"cvv"`
	ExpirationDate string `json:"expiration_date"`
	Password       string `json:"password"`
	Currency       string `json:"currency"`
	Selected       bool   `json:"selected"`
}

func GetWallets(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var wallets []models.Wallet
	database.DB.Where("user_id = ?", userID).Find(&wallets)
	for i := range wallets {
		wallets[i].ApplyMasks()
	}
	utils.OK(c, wallets)
}

func CreateWallet(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input WalletInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	validProviders := map[string]bool{"Nita": true, "MPesa": true, "Visa": true, "Mastercard": true, "PayPal": true}
	if !validProviders[input.Provider] {
		utils.BadRequest(c, "Invalid provider")
		return
	}

	currency := input.Currency
	if currency == "" {
		currency = "XOF"
	}

	wallet := models.Wallet{
		UserID:         userID,
		Provider:       input.Provider,
		Nature:         input.Nature,
		PhoneNumber:    input.PhoneNumber,
		Email:          input.Email,
		CardNumber:     input.CardNumber,
		CVV:            input.CVV,
		ExpirationDate: input.ExpirationDate,
		Currency:       currency,
		Balance:        0,
	}

	if input.Password != "" {
		hash, _ := utils.HashPassword(input.Password)
		wallet.Password = hash
	}

	if input.Selected {
		database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Update("selected", false)
		wallet.Selected = true
	}

	if err := database.DB.Create(&wallet).Error; err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "UNIQUE") {
			utils.BadRequest(c, "Wallet credentials already in use")
			return
		}
		utils.InternalError(c, "Failed to create wallet")
		return
	}

	wallet.ApplyMasks()
	utils.Created(c, wallet)
}

func UpdateWallet(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	walletID := c.Param("id")

	var wallet models.Wallet
	if err := database.DB.Where("id = ? AND user_id = ?", walletID, userID).First(&wallet).Error; err != nil {
		utils.NotFound(c, "Wallet not found")
		return
	}

	var input WalletInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if input.Nature != "" {
		updates["nature"] = input.Nature
	}
	if input.PhoneNumber != "" {
		updates["phone_number"] = input.PhoneNumber
	}
	if input.Email != "" {
		updates["email"] = input.Email
	}
	if input.Currency != "" {
		updates["currency"] = input.Currency
	}

	if input.Selected {
		database.DB.Model(&models.Wallet{}).Where("user_id = ? AND id != ?", userID, wallet.ID).Update("selected", false)
		updates["selected"] = true
	}

	database.DB.Model(&wallet).Updates(updates)
	wallet.ApplyMasks()
	utils.OK(c, wallet)
}

func DeleteWallet(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	walletID := c.Param("id")

	var wallet models.Wallet
	if err := database.DB.Where("id = ? AND user_id = ?", walletID, userID).First(&wallet).Error; err != nil {
		utils.NotFound(c, "Wallet not found")
		return
	}

	database.DB.Delete(&wallet)
	utils.OK(c, gin.H{"message": "Wallet deleted"})
}

func SelectWallet(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	walletID := c.Param("id")

	var wallet models.Wallet
	if err := database.DB.Where("id = ? AND user_id = ?", walletID, userID).First(&wallet).Error; err != nil {
		utils.NotFound(c, "Wallet not found")
		return
	}

	database.DB.Model(&models.Wallet{}).Where("user_id = ?", userID).Update("selected", false)
	database.DB.Model(&wallet).Update("selected", true)
	utils.OK(c, gin.H{"message": "Default wallet updated"})
}

func RefreshWalletBalance(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	walletID := c.Param("id")

	var wallet models.Wallet
	if err := database.DB.Where("id = ? AND user_id = ?", walletID, userID).First(&wallet).Error; err != nil {
		utils.NotFound(c, "Wallet not found")
		return
	}

	// Stub: mock balance from provider
	balances := map[string]float64{"Nita": 1000, "MPesa": 500, "Visa": 750, "Mastercard": 600, "PayPal": 300}
	balance := balances[wallet.Provider]

	database.DB.Model(&wallet).Update("balance", balance)
	utils.OK(c, gin.H{"balance": balance, "currency": wallet.Currency})
}
