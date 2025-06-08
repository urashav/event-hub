package http

import (
	"fmt"
	"github.com/urashav/event-hub/internal/dto/request"
	"github.com/urashav/event-hub/internal/dto/response"
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
		fmt.Println(err)
		httputils.ErrorResponse.Unauthorized(w, "Invalid credentials")
		return
	}
	httputils.SendSuccess(w, map[string]string{"token": token}, http.StatusOK)
}
