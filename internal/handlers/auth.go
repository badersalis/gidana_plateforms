package handlers

import (
	"regexp"
	"strings"
	"time"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	"github.com/gin-gonic/gin"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)
)

type RegisterInput struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password" binding:"required"`
}

type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Tous les champs obligatoires doivent être remplis")
		return
	}

	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.PhoneNumber = strings.TrimSpace(input.PhoneNumber)

	if len(input.FirstName) < 2 || len(input.FirstName) > 100 {
		utils.BadRequest(c, "Le prénom doit avoir entre 2 et 100 caractères")
		return
	}
	if len(input.LastName) < 2 || len(input.LastName) > 100 {
		utils.BadRequest(c, "Le nom doit avoir entre 2 et 100 caractères")
		return
	}
	if input.Email == "" && input.PhoneNumber == "" {
		utils.BadRequest(c, "Email ou numéro de téléphone requis")
		return
	}
	if input.Email != "" && !emailRegex.MatchString(input.Email) {
		utils.BadRequest(c, "Format d'email invalide")
		return
	}
	if input.PhoneNumber != "" && !phoneRegex.MatchString(input.PhoneNumber) {
		utils.BadRequest(c, "Format de téléphone invalide. Utilisez le format international (+XXXXXXXXXXX)")
		return
	}
	if len(input.Password) < 6 {
		utils.BadRequest(c, "Le mot de passe doit avoir au moins 6 caractères")
		return
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.InternalError(c, "Failed to hash password")
		return
	}

	user := models.User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PhoneNumber:  input.PhoneNumber,
		PasswordHash: hash,
		MemberSince:  time.Now(),
		Active:       true,
		Locale:       "fr",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "UNIQUE") {
			utils.BadRequest(c, "Cet email ou numéro de téléphone est déjà utilisé")
			return
		}
		utils.InternalError(c, "Failed to create account")
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email)
	utils.Created(c, gin.H{"user": user, "token": token})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "Email/téléphone et mot de passe requis")
		return
	}

	identifier := strings.TrimSpace(input.Identifier)
	if strings.Contains(identifier, "@") {
		identifier = strings.ToLower(identifier)
	}

	var user models.User
	if err := database.DB.Where("email = ? OR phone_number = ?", identifier, identifier).First(&user).Error; err != nil {
		utils.Unauthorized(c, "Identifiants incorrects")
		return
	}

	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		utils.Unauthorized(c, "Identifiants incorrects")
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Email)
	utils.OK(c, gin.H{"user": user, "token": token})
}

func GetMe(c *gin.Context) {
	user, _ := c.Get("user")
	utils.OK(c, user)
}
