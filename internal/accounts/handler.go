package accounts

import (
	"net/http"

	"github.com/simwawatim/book-api/internal/core"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var in RegisterInput
	if err := core.Decode(w, r, &in); err != nil {
		core.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.Register(r.Context(), in)
	if err != nil {
		core.WriteError(w, err)
		return
	}
	core.JSON(w, http.StatusCreated, res)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var in LoginInput
	if err := core.Decode(w, r, &in); err != nil {
		core.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.svc.Login(r.Context(), in)
	if err != nil {
		core.WriteError(w, err)
		return
	}
	core.JSON(w, http.StatusOK, res)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := core.UserID(r.Context())
	user, err := h.svc.Me(r.Context(), userID)
	if err != nil {
		core.WriteError(w, err)
		return
	}
	core.JSON(w, http.StatusOK, user)
}
