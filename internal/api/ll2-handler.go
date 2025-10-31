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

func (h *Handler) StartLL2AngecyUpdate(c *gin.Context) {
	err := h.ll2Server.UpdateAngecy(true)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "LL2 angecy update started"})
}

func (h *Handler) GetLL2Angecy(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	angecies, err := h.ll2Server.GetAngecyFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get angecies: "+err.Error())
		return
	}
	h.Json(c, angecies)
}

func (h *Handler) GetLL2LauncherFamilies(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	families, err := h.ll2Server.GetLauncherFamiliesFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get launcher families: "+err.Error())
		return
	}
	h.Json(c, families)
}

func (h *Handler) GetLL2Launchers(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	launchers, err := h.ll2Server.GetLaunchersFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get launchers: "+err.Error())
		return
	}
	h.Json(c, launchers)
}

func (h *Handler) StartLL2LauncherUpdate(c *gin.Context) {
	err := h.ll2Server.UpdateLaunchersAsync(true)
	if err != nil {
		h.Error(c, "start error:"+err.Error())
		return
	}
	h.Success(c, "ok")
}

func (h *Handler) StartLL2LauncherFamilyUpdate(c *gin.Context) {
	err := h.ll2Server.UpdateLauncherFamiliesAsync(true)
	if err != nil {
		h.Error(c, "start error:"+err.Error())
		return
	}
	h.Success(c, "ok")
}

func (h *Handler) StartLL2LocationUpdate(c *gin.Context) {
	err := h.ll2Server.UpdateLocationsAsync(true)
	if err != nil {
		h.Error(c, "start error:"+err.Error())
		return
	}
	h.Success(c, "ok")
}

func (h *Handler) GetLL2Locations(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	locations, err := h.ll2Server.GetLocationsFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get locations: "+err.Error())
		return
	}
	h.Json(c, locations)
}

func (h *Handler) GetLL2Pads(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	pads, err := h.ll2Server.GetPadsFromDB(limit, offset)
	if err != nil {
		h.Error(c, "failed to get pads: "+err.Error())
		return
	}
	h.Json(c, pads)
}

func (h *Handler) StartLL2PadUpdate(c *gin.Context) {
	err := h.ll2Server.UpdatePadsAsync(true)
	if err != nil {
		h.Error(c, "start error:"+err.Error())
		return
	}
	h.Success(c, "ok")
}
