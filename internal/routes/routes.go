package routes

import (
	"net/http"
	"strings"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/badersalis/gidana_backend/internal/handlers"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	allowedOrigins := make(map[string]bool)
	for _, o := range strings.Split(config.App.AllowedOrigins, ",") {
		allowedOrigins[strings.TrimSpace(o)] = true
	}

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// React Native native (iOS/Android) sends no origin — always allow.
			if origin == "" || origin == "null" {
				return true
			}
			return allowedOrigins[origin] || allowedOrigins["*"]
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	r.Static("/uploads", "./uploads")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "1.0.0"})
	})

	api := r.Group("/api/v1")

	// Auth
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.GET("/me", middleware.Auth(), handlers.GetMe)
	}

	// WebSocket
	r.GET("/ws", handlers.ServeWS)

	// User
	users := api.Group("/users", middleware.Auth())
	{
		users.PUT("/profile", handlers.UpdateProfile)
		users.POST("/profile-picture", handlers.UploadProfilePicture)
		users.PUT("/password", handlers.ChangePassword)
		users.PATCH("/push-token", handlers.UpdatePushToken)
		users.DELETE("/profile", handlers.RequestDeleteAccount)
	}

	// Properties
	props := api.Group("/properties")
	{
		props.GET("", middleware.OptionalAuth(), handlers.ListProperties)
		props.GET("/featured", handlers.GetFeaturedProperties)
		props.GET("/:id", middleware.OptionalAuth(), handlers.GetProperty)
		props.POST("", middleware.Auth(), handlers.CreateProperty)
		props.PUT("/:id", middleware.Auth(), handlers.UpdateProperty)
		props.DELETE("/:id", middleware.Auth(), handlers.DeleteProperty)
		props.PATCH("/:id/availability", middleware.Auth(), handlers.ToggleAvailability)
		props.GET("/my/listings", middleware.Auth(), handlers.MyProperties)

		// Property images
		props.POST("/:id/images", middleware.Auth(), handlers.AddPropertyImage)

		// Reviews
		props.GET("/:id/reviews", handlers.GetPropertyReviews)
		props.POST("/:id/reviews", middleware.Auth(), handlers.CreateReview)
	}

	// Images
	images := api.Group("/images", middleware.Auth())
	{
		images.DELETE("/:id", handlers.DeletePropertyImage)
		images.PATCH("/:id/main", handlers.SetMainImage)
	}

	// Reviews
	reviews := api.Group("/reviews", middleware.Auth())
	{
		reviews.DELETE("/:id", handlers.DeleteReview)
	}

	// Favorites
	favs := api.Group("/favorites", middleware.Auth())
	{
		favs.GET("", handlers.GetFavorites)
		favs.POST("/:id/toggle", handlers.ToggleFavorite)
	}

	// Rentals
	rentals := api.Group("/rentals", middleware.Auth())
	{
		rentals.GET("", handlers.GetMyRentals)
		rentals.POST("", handlers.CreateRental)
		rentals.PATCH("/:id/status", handlers.UpdateRentalStatus)
	}

	// Wallets
	wallets := api.Group("/wallets", middleware.Auth())
	{
		wallets.GET("", handlers.GetWallets)
		wallets.POST("", handlers.CreateWallet)
		wallets.PUT("/:id", handlers.UpdateWallet)
		wallets.DELETE("/:id", handlers.DeleteWallet)
		wallets.PATCH("/:id/select", handlers.SelectWallet)
		wallets.POST("/:id/refresh-balance", handlers.RefreshWalletBalance)
	}

	// Transactions & Payments
	txs := api.Group("/transactions", middleware.Auth())
	{
		txs.GET("", handlers.GetTransactions)
		txs.POST("/pay-service", handlers.PayService)
		txs.POST("/transfer", handlers.TransferMoney)
	}

	// Alerts
	alerts := api.Group("/alerts", middleware.Auth())
	{
		alerts.GET("", handlers.GetAlerts)
		alerts.POST("", handlers.CreateAlert)
		alerts.PUT("/:id", handlers.UpdateAlert)
		alerts.DELETE("/:id", handlers.DeleteAlert)
	}

	// Messaging
	convs := api.Group("/conversations", middleware.Auth())
	{
		convs.GET("", handlers.GetConversations)
		convs.POST("", handlers.StartConversation)
		convs.GET("/:id", handlers.GetConversation)
		convs.POST("/:id/messages", handlers.SendMessage)
		convs.DELETE("/:id/messages/:msgID", handlers.DeleteMessage)
	}

	// Search
	search := api.Group("/search")
	{
		search.GET("/suggestions", handlers.GetSearchSuggestions)
		search.POST("/history", middleware.OptionalAuth(), handlers.SaveSearchHistory)
		search.GET("/history", middleware.Auth(), handlers.GetSearchHistory)
		search.DELETE("/history", middleware.Auth(), handlers.DeleteSearchHistory)
	}
}
