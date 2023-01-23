package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hugosrc/shortlink/internal/core/port"
	"github.com/hugosrc/shortlink/internal/util"
)

type LinkHandler struct {
	auth port.Auth
	svc  port.LinkService
}

func NewLinkHandler(auth port.Auth, svc port.LinkService) *LinkHandler {
	return &LinkHandler{
		auth: auth,
		svc:  svc,
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
		handleError(w, err, "An internal error has occurred. Please try again later.")
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url"`
}

func (h *LinkHandler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := h.auth.Authenticate(r, w)
	if err != nil {
		handleError(w, err, "Invalid authentication credentials")
		return
	}

	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, util.WrapErrorf(err, util.ErrCodeInvalidArgument, "json decode"),
			"Invalid request format")
		return
	}

	link, err := h.svc.Create(r.Context(), req.OriginalURL, userID)
	if err != nil {
		handleError(w, err, "An internal error has occurred. Please try again later.")
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&link)
}

type UpdateLinkRequest struct {
	OriginalURL string `json:"original_url"`
}

func (h *LinkHandler) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := h.auth.Authenticate(r, w)
	if err != nil {
		handleError(w, err, "Invalid authentication credentials")
		return
	}

	var req UpdateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, util.WrapErrorf(err, util.ErrCodeInvalidArgument, "json decode"),
			"Invalid request format")
		return
	}

	vars := mux.Vars(r)
	link, err := h.svc.Update(r.Context(), vars["hash"], req.OriginalURL, userID)
	if err != nil {
		handleError(w, err, "An internal error has occurred. Please try again later.")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&link)
}

func (h *LinkHandler) delete(w http.ResponseWriter, r *http.Request) {
	userID, err := h.auth.Authenticate(r, w)
	if err != nil {
		handleError(w, err, "Invalid authentication credentials")
		return
	}

	vars := mux.Vars(r)

	err = h.svc.Delete(r.Context(), vars["hash"], userID)
	if err != nil {
		handleError(w, err, "An internal error has occurred. Please try again later.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
