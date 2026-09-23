package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/middleware"
	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input domain.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.BadRequest(w, "invalid request body", err.Error())
		return
	}

	if input.GetIdentifier() == "" || input.Password == "" {
		response.BadRequest(w, "username and password are required", nil)
		return
	}

	authResp, err := h.authService.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Unauthorized(w, "invalid username or password")
			return
		}
		response.InternalServerError(w, "login failed")
		return
	}

	response.OK(w, "authenticated successfully", authResp)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "user not found")
			return
		}
		response.InternalServerError(w, "failed to get current user")
		return
	}

	response.OK(w, "current user retrieved", user)
}
