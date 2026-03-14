package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/utils"
)

// TaskHandler holds dependencies for task HTTP handlers.
type TaskHandler struct {
	taskDAO *dao.TaskDAO
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(taskDAO *dao.TaskDAO) *TaskHandler {
	return &TaskHandler{taskDAO: taskDAO}
}

// createTaskRequest is the expected body for POST /tasks.
type createTaskRequest struct {
	Title string `json:"title" binding:"required"`
}

// CreateTask handles POST /tasks – create a new task.
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	task := &dao.Task{Title: req.Title}
	if err := h.taskDAO.CreateTask(task); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create task")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "task created", task)
}

// GetTask handles GET /tasks/:id – retrieve a task by ID.
func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := h.taskDAO.GetTask(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "task not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", task)
}

// updateTaskRequest is the expected body for PUT /tasks/:id.
type updateTaskRequest struct {
	Title string `json:"title"`
	Done  *bool  `json:"done"`
}

// UpdateTask handles PUT /tasks/:id – update a task.
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := h.taskDAO.GetTask(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "task not found")
		return
	}

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Done != nil {
		task.Done = *req.Done
	}

	if err := h.taskDAO.UpdateTask(task); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "task updated", task)
}

// DeleteTask handles DELETE /tasks/:id – remove a task.
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid task id")
		return
	}

	if err := h.taskDAO.DeleteTask(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to delete task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "task deleted", nil)
}