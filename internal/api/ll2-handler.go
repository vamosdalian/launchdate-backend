package api

import "github.com/gin-gonic/gin"

func (h *Handler) GetLL2Launches(c *gin.Context) {
	// Implementation for fetching LL2 launches goes here
}

func (h *Handler) StartLL2LaunchUpdate(c *gin.Context) {
	err := h.ll2Server.UpdateLaunches(true)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "LL2 launch update started"})
}
