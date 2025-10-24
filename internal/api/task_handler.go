package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

// CreateTask creates a new task
// @Summary Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.CreateTaskRequest true "Task data"
// @Success 201 {object} models.Task
// @Router /api/v1/tasks [post]
func (h *Handler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.CreateTask(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("failed to create task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTask retrieves a task by ID
// @Summary Get a task
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Router /api/v1/tasks/{id} [get]
func (h *Handler) GetTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	task, err := h.taskService.GetTask(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("failed to get task")
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// ListLaunchTasks retrieves all tasks for a launch
// @Summary List tasks for a launch
// @Tags tasks
// @Produce json
// @Param launch_id path int true "Launch ID"
// @Success 200 {array} models.Task
// @Router /api/v1/launches/{launch_id}/tasks [get]
func (h *Handler) ListLaunchTasks(c *gin.Context) {
	launchID, err := strconv.ParseInt(c.Param("launch_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid launch_id"})
		return
	}

	tasks, err := h.taskService.ListTasksByLaunchID(c.Request.Context(), launchID)
	if err != nil {
		h.logger.WithError(err).Error("failed to list tasks")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// UpdateTask updates a task
// @Summary Update a task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.UpdateTaskRequest true "Task data"
// @Success 200 {object} gin.H
// @Router /api/v1/tasks/{id} [put]
func (h *Handler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.taskService.UpdateTask(c.Request.Context(), id, &req); err != nil {
		h.logger.WithError(err).Error("failed to update task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

// DeleteTask deletes a task
// @Summary Delete a task
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204
// @Router /api/v1/tasks/{id} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.taskService.DeleteTask(c.Request.Context(), id); err != nil {
		h.logger.WithError(err).Error("failed to delete task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	c.Status(http.StatusNoContent)
}
