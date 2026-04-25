package main

import (
	"log"
	"os"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/routes"
	"github.com/badersalis/gidana_backend/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	if config.App.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	database.Connect()

	if config.App.UseFirebase {
		if err := storage.Init(); err != nil {
			log.Fatalf("Firebase Storage init failed: %v", err)
		}
		log.Println("Firebase Storage initialized")
	}

	if err := os.MkdirAll(config.App.UploadDir, 0755); err != nil {
		log.Printf("Warning: could not create upload dir: %v", err)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	routes.Setup(r)

	port := config.App.Port
	log.Printf("Gidana API server starting on port %s (env: %s)", port, config.App.AppEnv)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
