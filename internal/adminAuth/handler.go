package adminauth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/varnit-ta/PlacementLog/pkg/jwt"
	"github.com/varnit-ta/PlacementLog/pkg/utils"
)

type AdminAuthHandler struct {
	service *AdminService
}

func NewAdminAuthHandler(service *AdminService) *AdminAuthHandler {
	return &AdminAuthHandler{service: service}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type responsePayload struct {
	UserID string `json:"username"`
	Token  string `json:"token"`
}

func (h AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, admin, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := responsePayload{
		UserID: admin.ID,
		Token:  token,
	}

	utils.WriteJSON(w, resp, http.StatusOK)
}

func (h AdminAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "unauthorized: missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	_, err := jwt.ValidateJwtToken(token)
	if err != nil {
		http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tokenStr, admin, err := h.service.Register(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := responsePayload{
		UserID: admin.ID,
		Token:  tokenStr,
	}

	utils.WriteJSON(w, resp, http.StatusCreated)
}
