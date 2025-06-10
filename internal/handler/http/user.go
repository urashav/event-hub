package http

import (
	"encoding/json"
	"github.com/urashav/event-hub/internal/dto/request"
	"github.com/urashav/event-hub/internal/dto/response"
	"github.com/urashav/event-hub/internal/models"
	"github.com/urashav/event-hub/internal/service"
	httputils "github.com/urashav/event-hub/pkg/httputilst"
	"net/http"
)

type UserHandler struct {
	service service.UsersService
}

func NewUserHandler(service *service.UsersService) *UserHandler {
	return &UserHandler{
		service: *service,
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputils.ErrorResponse.MethodNotAllowed(w, "Method not allowed")
		return
	}

	var req request.SignUpRequest

	if err := httputils.DecodeJSON(r, &req); err != nil {
		httputils.ErrorResponse.BadRequest(w, "Invalid JSON")
		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		httputils.ErrorResponse.BadRequest(w, "Validation Error: "+err.Error())
		return
	}

	user := req.ToModel()

	createdUser, err := h.service.CreateUser(r.Context(), *user)
	if err != nil {
		httputils.ErrorResponse.InternalError(w, "Error creating user: "+err.Error())
		return
	}

	resp := response.FromModel(createdUser)
	httputils.SendSuccess(w, resp, http.StatusCreated)
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req request.SignInRequest
	if err := httputils.DecodeJSON(r, &req); err != nil {
		httputils.ErrorResponse.BadRequest(w, "Invalid JSON")
		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		httputils.ErrorResponse.BadRequest(w, "Validation Error: "+err.Error())
	}

	token, err := h.service.AuthenticateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		httputils.ErrorResponse.Unauthorized(w, "Invalid credentials")
		return
	}
	httputils.SendSuccess(w, map[string]string{"token": token}, http.StatusOK)
}
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputils.ErrorResponse.MethodNotAllowed(w, "Method not allowed")
		return
	}

	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		httputils.ErrorResponse.InternalError(w, "Failed to get users")
		return
	}

	httputils.SendSuccess(w, users, http.StatusOK)
}

func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httputils.ErrorResponse.MethodNotAllowed(w, "Method not allowed")
		return
	}

	var req struct {
		UserID int         `json:"user_id"`
		Role   models.Role `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorResponse.BadRequest(w, "Invalid request body")
		return
	}

	adminID := r.Context().Value("user_id").(int)
	err := h.service.UpdateUserRole(r.Context(), adminID, req.UserID, req.Role)
	if err != nil {
		httputils.ErrorResponse.InternalError(w, "Failed to update user role")
		return
	}

	httputils.SendSuccess(w, map[string]string{"message": "Role updated successfully"}, http.StatusOK)
}
