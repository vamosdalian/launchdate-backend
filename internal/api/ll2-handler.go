package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetLL2Launches(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	launches, err := h.ll2Server.GetLaunchesFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get launches: "+err.Error())
		return
	}
	h.Json(c, launches)
}

func (h *Handler) StartLL2LaunchUpdate(c *gin.Context) {
	err := h.ll2Server.UpdateLaunches(true)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "LL2 launch update started"})
}
