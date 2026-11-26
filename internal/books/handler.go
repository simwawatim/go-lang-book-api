package books

import (
	"net/http"
	"strconv"

	"github.com/simwawatim/book-api/internal/core"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	res, err := h.svc.List(r.Context(), ListParams{
		Search: q.Get("search"),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		core.WriteError(w, err)
		return
	}
	core.JSON(w, http.StatusOK, res)
}

// GET /api/books/{id}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	book, err := h.svc.Get(r.Context(), id)
	if err != nil {
		core.WriteError(w, err)
		return
	}
	core.JSON(w, http.StatusOK, book)
}

// POST /api/books  (requires auth)
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := core.UserID(r.Context())

	var in BookInput
	if err := core.Decode(w, r, &in); err != nil {
		core.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	book, err := h.svc.Create(r.Context(), userID, in)
	if err != nil {
		core.WriteError(w, err)
		return
	}
	core.JSON(w, http.StatusCreated, book)
}

// PUT /api/books/{id}  (requires auth, owner only)
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := core.UserID(r.Context())
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var in BookInput
	if err := core.Decode(w, r, &in); err != nil {
		core.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	book, err := h.svc.Update(r.Context(), userID, id, in)
	if err != nil {
		core.WriteError(w, err)
		return
	}
	core.JSON(w, http.StatusOK, book)
}

// DELETE /api/books/{id}  (requires auth, owner only)
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := core.UserID(r.Context())
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := h.svc.Delete(r.Context(), userID, id); err != nil {
		core.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Error(w, http.StatusBadRequest, "invalid book id")
		return 0, false
	}
	return id, true
}
