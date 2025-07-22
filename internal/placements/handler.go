package placements

import (
	"net/http"

	"github.com/varnit-ta/PlacementLog/pkg/utils"
)

type PlacementsHandler struct {
	srv *PlacementsService
}

func NewPlacementsHandler(srv *PlacementsService) *PlacementsHandler {
	return &PlacementsHandler{srv: srv}
}

type PlacementRequest struct {
	Company       string   `json:"company"`
	CTC           float64  `json:"ctc"`
	PlacementDate string   `json:"placement_date"`
	Students      []string `json:"students"`
}

type PlacementResponse struct {
	PlacementID   int           `json:"placement_id"`
	Company       string        `json:"company"`
	CTC           float64       `json:"ctc"`
	PlacementDate string        `json:"placement_date"`
	BranchCounts  []BranchCount `json:"branch_counts"`
}

// POST /placements (admin only, enforced by router middleware)
func (h *PlacementsHandler) AddPlacement(w http.ResponseWriter, r *http.Request) {
	var req PlacementRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteError(w, err)
		return
	}
	resp, err := h.srv.AddPlacement(req)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	utils.WriteJSON(w, resp, http.StatusCreated)
}

// GET /placements (all users)
func (h *PlacementsHandler) GetAllPlacements(w http.ResponseWriter, r *http.Request) {
	placementsList, err := h.srv.GetAllPlacements()
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	utils.WriteJSON(w, placementsList, http.StatusOK)
}
