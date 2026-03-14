package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-api-golang/config"
	"mini-api-golang/internal/service"
	"mini-api-golang/internal/utils"
)

// UserHandler holds dependencies for user HTTP handlers.
type UserHandler struct {
	userService *service.UserService
	jwtSecret   string
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtSecret:   cfg.JWT.Secret,
	}
}

// registerRequest is the expected body for POST /register.
type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// loginRequest is the expected body for POST /login.
type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles POST /register – create a new user account.
func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusConflict, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "user registered", user)
}

// Login handles POST /login and POST /api/auth/login – authenticate and return a JWT.
func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := utils.GenerateJWT(h.jwtSecret, strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "login successful", gin.H{
		"token": token,
		"user":  user,
	})
}

// Me handles GET /api/me – return the currently authenticated user.
func (h *UserHandler) Me(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	user, err := h.userService.GetByID(uint(userID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", user)
}

// GetUser handles GET /users/:id – retrieve a user by ID.
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", user)
}

// updateUserRequest is the expected body for PUT /users/:id.
type updateUserRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
}

// UpdateUser handles PUT /users/:id – update user fields.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if err := h.userService.Update(user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update user")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user updated", user)
}

// DeleteUser handles DELETE /users/:id – remove a user.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.userService.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to delete user")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user deleted", nil)
}

// parseUintParam extracts and parses a named URL parameter as uint.
func parseUintParam(c *gin.Context, name string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}