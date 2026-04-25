package handlers

import (
	"strconv"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

var servicePlans = map[string]map[string]float64{
	"starlink": {"Basic": 50, "Standard": 150, "Premium": 500},
	"canal+":   {"Standard": 30},
}

type PayServiceInput struct {
	Service         string  `json:"service" binding:"required"`
	ServiceProvider string  `json:"service_provider" binding:"required"`
	Plan            string  `json:"plan"`
	WalletID        uint    `json:"wallet_id" binding:"required"`
	Amount          float64 `json:"amount"`
}

type TransferInput struct {
	WalletID    uint    `json:"wallet_id" binding:"required"`
	Recipient   string  `json:"recipient" binding:"required"` // phone or email
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Provider    string  `json:"provider" binding:"required"`
}

func GetTransactions(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.Transaction{}).Where("user_id = ?", userID).Count(&total)

	var txs []models.Transaction
	database.DB.Where("user_id = ?", userID).
		Order("created_at desc").
		Offset(offset).Limit(pageSize).
		Preload("Wallet").
		Find(&txs)

	utils.Paginated(c, txs, total, page, pageSize)
}

func PayService(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input PayServiceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var wallet models.Wallet
	if err := database.DB.Where("id = ? AND user_id = ?", input.WalletID, userID).First(&wallet).Error; err != nil {
		utils.NotFound(c, "Wallet not found")
		return
	}

	amount := input.Amount
	if amount == 0 {
		plans, ok := servicePlans[input.ServiceProvider]
		if !ok {
			utils.BadRequest(c, "Unknown service provider")
			return
		}
		base, ok := plans[input.Plan]
		if !ok {
			utils.BadRequest(c, "Unknown plan")
			return
		}
		amount = base * 1.1 // 10% service fee
	}

	if wallet.Balance < amount {
		utils.BadRequest(c, "Insufficient balance")
		return
	}

	database.DB.Model(&wallet).Update("balance", wallet.Balance-amount)

	tx := models.Transaction{
		UserID:          userID,
		WalletID:        wallet.ID,
		Amount:          amount,
		Nature:          "expense",
		Service:         input.Service,
		ServiceProvider: input.ServiceProvider,
		Currency:        wallet.Currency,
		Status:          models.StatusDone,
	}
	database.DB.Create(&tx)

	utils.OK(c, gin.H{
		"message":         "Payment successful",
		"amount":          amount,
		"currency":        wallet.Currency,
		"new_balance":     wallet.Balance - amount,
		"transaction_id":  tx.ID,
	})
}

func TransferMoney(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input TransferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var senderWallet models.Wallet
	if err := database.DB.Where("id = ? AND user_id = ?", input.WalletID, userID).First(&senderWallet).Error; err != nil {
		utils.NotFound(c, "Wallet not found")
		return
	}

	// Prevent self-transfer
	var recipientWallet models.Wallet
	if input.Provider == "Nita" || input.Provider == "MPesa" {
		if senderWallet.PhoneNumber == input.Recipient {
			utils.BadRequest(c, "Cannot transfer to yourself")
			return
		}
		if err := database.DB.Where("phone_number = ? AND provider = ?", input.Recipient, input.Provider).First(&recipientWallet).Error; err != nil {
			utils.NotFound(c, "Recipient wallet not found")
			return
		}
	} else if input.Provider == "PayPal" {
		if senderWallet.Email == input.Recipient {
			utils.BadRequest(c, "Cannot transfer to yourself")
			return
		}
		if err := database.DB.Where("email = ? AND provider = 'PayPal'", input.Recipient).First(&recipientWallet).Error; err != nil {
			utils.NotFound(c, "Recipient wallet not found")
			return
		}
	}

	if senderWallet.Currency != recipientWallet.Currency {
		utils.BadRequest(c, "Currency mismatch between wallets")
		return
	}

	if senderWallet.Balance < input.Amount {
		utils.BadRequest(c, "Insufficient balance")
		return
	}

	// Debit sender
	database.DB.Model(&senderWallet).Update("balance", senderWallet.Balance-input.Amount)

	outTx := models.Transaction{
		UserID:          userID,
		WalletID:        senderWallet.ID,
		Amount:          input.Amount,
		Nature:          "expense",
		Service:         "transfer",
		ServiceProvider: input.Provider,
		Currency:        senderWallet.Currency,
		Status:          models.StatusDone,
	}
	database.DB.Create(&outTx)

	// Credit recipient
	database.DB.Model(&recipientWallet).Update("balance", recipientWallet.Balance+input.Amount)

	inTx := models.Transaction{
		UserID:               recipientWallet.UserID,
		WalletID:             recipientWallet.ID,
		Amount:               input.Amount,
		Nature:               "income",
		Service:              "transfer",
		ServiceProvider:      input.Provider,
		Currency:             recipientWallet.Currency,
		Status:               models.StatusDone,
		RelatedTransactionID: &outTx.ID,
	}
	database.DB.Create(&inTx)

	database.DB.Model(&outTx).Update("related_transaction_id", inTx.ID)

	utils.OK(c, gin.H{
		"message":        "Transfer successful",
		"amount":         input.Amount,
		"currency":       senderWallet.Currency,
		"new_balance":    senderWallet.Balance - input.Amount,
		"transaction_id": outTx.ID,
	})
}
