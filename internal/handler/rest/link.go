package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hugosrc/shortlink/internal/core/port"
)

type LinkHandler struct {
	svc port.LinkService
}

func NewLinkHandler(svc port.LinkService) *LinkHandler {
	return &LinkHandler{
		svc: svc,
	}
}

func (h *LinkHandler) Register(r *mux.Router) {
	r.HandleFunc("/{hash}", h.show).Methods(http.MethodGet)
	r.HandleFunc("/api/shortlink", h.create).Methods(http.MethodPost)
	r.HandleFunc("/api/shortlink/{hash}", h.update).Methods(http.MethodPut)
	r.HandleFunc("/api/shortlink/{hash}", h.delete).Methods(http.MethodDelete)
}

func (h *LinkHandler) show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	url, err := h.svc.FindByHash(r.Context(), vars["hash"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url"`
}

func (h *LinkHandler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link, err := h.svc.Create(r.Context(), req.OriginalURL, "803681d8-7d3f-4078-989a-ae535a073624") // TODO: Add Auth System
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(&link); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type UpdateLinkRequest struct {
	OriginalURL string `json:"original_url"`
}

func (h *LinkHandler) update(w http.ResponseWriter, r *http.Request) {
	var req UpdateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	link, err := h.svc.Update(r.Context(), vars["hash"], req.OriginalURL, "803681d8-7d3f-4078-989a-ae535a073624")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(&link); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *LinkHandler) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := h.svc.Delete(r.Context(), vars["hash"], "803681d8-7d3f-4078-989a-ae535a073624")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
