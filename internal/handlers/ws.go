package handlers

import (
	"github.com/badersalis/gidana_backend/internal/utils"
	appws "github.com/badersalis/gidana_backend/internal/ws"
	"github.com/gin-gonic/gin"
)

// ServeWS upgrades an HTTP connection to WebSocket.
// Auth uses ?token= because browsers cannot send Authorization headers in WS handshakes.
func ServeWS(c *gin.Context) {
	tokenStr := c.Query("token")
	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		utils.Unauthorized(c, "Invalid token")
		return
	}

	conn, err := appws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	appws.H.Connect(claims.UserID, conn)
	defer appws.H.Disconnect(claims.UserID)

	// Read pump: discard client messages but keep the connection alive for pings.
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
