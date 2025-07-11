package userauth

import (
	"net/http"

	"github.com/varnit-ta/PlacementLog/pkg/utils"
)

type UserAuthHandler struct {
	srv *UserAuthService
}

type requestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type responsePayload struct {
	UserID string `json:"username"`
	Token  string `json:"token"`
}

func NewUserAuthHandler(srv *UserAuthService) *UserAuthHandler {
	return &UserAuthHandler{
		srv: srv,
	}
}

func (h *UserAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload requestPayload

	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.WriteError(w, err)
		return
	}

	token, userId, err := h.srv.Login(payload.Username, payload.Password)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	resp := responsePayload{
		UserID: userId,
		Token:  token,
	}

	utils.WriteJSON(w, resp, http.StatusOK)
}

func (h *UserAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload requestPayload

	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.WriteError(w, err)
		return
	}

	token, userId, err := h.srv.Register(payload.Username, payload.Password)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	resp := responsePayload{
		UserID: userId,
		Token:  token,
	}

	utils.WriteJSON(w, resp, http.StatusCreated)
}
