package http

import (
	"delivery-service/internal/domain"
	"delivery-service/internal/usecase"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CourierHandler struct {
	uc *usecase.CourierUseCase
}

func NewCourierHandler(uc *usecase.CourierUseCase) *CourierHandler {
	return &CourierHandler{uc: uc}
}

func (h *CourierHandler) Register(r chi.Router) {
	r.Route("/couriers", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/", h.list)
		r.Get("/{id}", h.getByID)
		r.Patch("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

func (h *CourierHandler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateCourierInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	courier, err := h.uc.Create(input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, courier)
}

func (h *CourierHandler) list(w http.ResponseWriter, r *http.Request) {
	couriers, err := h.uc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, couriers)
}

func (h *CourierHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	courier, err := h.uc.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrCourierNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, courier)
}

func (h *CourierHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var input domain.UpdateCourierInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	courier, err := h.uc.Update(id, input)
	if err != nil {
		if errors.Is(err, domain.ErrCourierNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, courier)
}

func (h *CourierHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.uc.Delete(id); err != nil {
		if errors.Is(err, domain.ErrCourierNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
