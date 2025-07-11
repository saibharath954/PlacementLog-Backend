package posts

import (
	"net/http"

	"github.com/varnit-ta/PlacementLog/pkg/utils"
)

type PostsHandler struct {
	srv *PostsService
}

func NewPostsHandler(srv *PostsService) *PostsHandler {
	return &PostsHandler{
		srv: srv,
	}
}

type createPostRequest struct {
	PostBody map[string]any `json:"post_body"`
}

func (h *PostsHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	var req createPostRequest

	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteError(w, err)
		return
	}

	userId := r.Header.Get("X-User-ID")

	if userId == "" {
		utils.WriteError(w, http.ErrNoCookie)
		return
	}

	post, err := h.srv.AddPost(userId, req.PostBody)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, post, http.StatusCreated)
}

func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	postId := query.Get("id")

	if postId == "" {
		utils.WriteError(w, http.ErrMissingFile)
		return
	}

	err := h.srv.DeletePost(postId)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "post deleted"}, http.StatusOK)
}

func (h *PostsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.srv.GetAll()
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, posts, http.StatusOK)
}

func (h *PostsHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userId := query.Get("user_id")

	if userId == "" {
		utils.WriteError(w, http.ErrMissingFile)
		return
	}

	posts, err := h.srv.GetByUser(userId)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, posts, http.StatusOK)
}
