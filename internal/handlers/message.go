package handlers

import (
	"fmt"
	"time"

	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/models"
	"github.com/badersalis/gidana_backend/internal/utils"
	appws "github.com/badersalis/gidana_backend/internal/ws"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StartConversationInput struct {
	RecipientID uint   `json:"recipient_id" binding:"required"`
	PropertyID  *uint  `json:"property_id"`
	Message     string `json:"message" binding:"required,min=1"`
}

type SendMessageInput struct {
	Content string `json:"content" binding:"required,min=1"`
}

// notifyRecipient delivers a new-message event via WebSocket and, when the
// recipient is not connected, falls back to an Expo push notification.
func notifyRecipient(recipientID uint, senderName string, msg models.Message) {
	appws.H.Emit(recipientID, appws.Event{Type: "new_message", Data: msg})

	if !appws.H.IsOnline(recipientID) {
		var recipient models.User
		if err := database.DB.First(&recipient, recipientID).Error; err == nil && recipient.ExpoPushToken != "" {
			utils.SendExpoPush(
				recipient.ExpoPushToken,
				senderName,
				msg.Content,
				map[string]any{"conversation_id": msg.ConversationID},
			)
		}
	}
}

// POST /conversations
func StartConversation(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var input StartConversationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if input.RecipientID == userID {
		utils.BadRequest(c, "Cannot start a conversation with yourself")
		return
	}

	var recipient models.User
	if err := database.DB.First(&recipient, input.RecipientID).Error; err != nil {
		utils.NotFound(c, "Recipient not found")
		return
	}

	if input.PropertyID != nil {
		var prop models.Property
		if err := database.DB.First(&prop, *input.PropertyID).Error; err != nil {
			utils.NotFound(c, "Property not found")
			return
		}
	}

	// Find existing conversation between these two users for this property (either direction)
	var conv models.Conversation
	q := database.DB.Where(
		"(initiator_id = ? AND recipient_id = ?) OR (initiator_id = ? AND recipient_id = ?)",
		userID, input.RecipientID, input.RecipientID, userID,
	)
	if input.PropertyID != nil {
		q = q.Where("property_id = ?", *input.PropertyID)
	} else {
		q = q.Where("property_id IS NULL")
	}

	if err := q.First(&conv).Error; err != nil {
		conv = models.Conversation{
			InitiatorID: userID,
			RecipientID: input.RecipientID,
			PropertyID:  input.PropertyID,
		}
		database.DB.Create(&conv)
	}

	msg := models.Message{
		ConversationID: conv.ID,
		SenderID:       userID,
		Content:        input.Message,
	}
	database.DB.Create(&msg)
	database.DB.Model(&models.Conversation{}).Where("id = ?", conv.ID).UpdateColumn("updated_at", time.Now())
	database.DB.Preload("Sender").First(&msg, msg.ID)

	database.DB.
		Preload("Initiator").
		Preload("Recipient").
		Preload("Property").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Sender").Order("messages.created_at ASC")
		}).
		First(&conv, conv.ID)

	senderName := fmt.Sprintf("%s %s", msg.Sender.FirstName, msg.Sender.LastName)
	notifyRecipient(input.RecipientID, senderName, msg)

	utils.Created(c, conv)
}

// GET /conversations
func GetConversations(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var convs []models.Conversation
	database.DB.
		Where("initiator_id = ? OR recipient_id = ?", userID, userID).
		Preload("Initiator").
		Preload("Recipient").
		Preload("Property").
		Order("updated_at DESC").
		Find(&convs)

	for i := range convs {
		var lastMsg models.Message
		if err := database.DB.Where("conversation_id = ?", convs[i].ID).
			Preload("Sender").
			Order("created_at DESC").
			First(&lastMsg).Error; err == nil {
			convs[i].LastMessage = &lastMsg
		}

		var unread int64
		database.DB.Model(&models.Message{}).
			Where("conversation_id = ? AND sender_id != ? AND is_read = false", convs[i].ID, userID).
			Count(&unread)
		convs[i].UnreadCount = int(unread)
	}

	utils.OK(c, convs)
}

// GET /conversations/:id
func GetConversation(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	convID := c.Param("id")

	var conv models.Conversation
	if err := database.DB.
		Preload("Initiator").
		Preload("Recipient").
		Preload("Property").
		First(&conv, convID).Error; err != nil {
		utils.NotFound(c, "Conversation not found")
		return
	}

	if conv.InitiatorID != userID && conv.RecipientID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	var msgs []models.Message
	database.DB.Where("conversation_id = ?", conv.ID).
		Preload("Sender").
		Order("created_at ASC").
		Find(&msgs)
	conv.Messages = msgs

	// Mark incoming messages as read
	database.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = false", conv.ID, userID).
		Update("is_read", true)

	utils.OK(c, conv)
}

// POST /conversations/:id/messages
func SendMessage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	convID := c.Param("id")

	var conv models.Conversation
	if err := database.DB.First(&conv, convID).Error; err != nil {
		utils.NotFound(c, "Conversation not found")
		return
	}

	if conv.InitiatorID != userID && conv.RecipientID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	var input SendMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	msg := models.Message{
		ConversationID: conv.ID,
		SenderID:       userID,
		Content:        input.Content,
	}
	database.DB.Create(&msg)
	database.DB.Model(&models.Conversation{}).Where("id = ?", conv.ID).UpdateColumn("updated_at", time.Now())
	database.DB.Preload("Sender").First(&msg, msg.ID)

	recipientID := conv.RecipientID
	if userID == conv.RecipientID {
		recipientID = conv.InitiatorID
	}
	senderName := fmt.Sprintf("%s %s", msg.Sender.FirstName, msg.Sender.LastName)
	notifyRecipient(recipientID, senderName, msg)

	utils.Created(c, msg)
}

// DELETE /conversations/:id/messages/:msgID
func DeleteMessage(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	msgID := c.Param("msgID")

	var msg models.Message
	if err := database.DB.First(&msg, msgID).Error; err != nil {
		utils.NotFound(c, "Message not found")
		return
	}

	if msg.SenderID != userID {
		utils.Forbidden(c, "Not authorized")
		return
	}

	database.DB.Delete(&msg)
	utils.OK(c, gin.H{"message": "Message deleted"})
}
