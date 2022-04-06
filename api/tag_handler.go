package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/ddddddO/memo/api/adapter"
)

type tagUsecase interface {
	List(userID int) ([]adapter.Tag, error)
	Detail(tagID int) (*adapter.Tag, error)
	Update(updatedTag adapter.Tag) error
	Delete(tagID int) error
	Create(createTag adapter.Tag) error
}

type tagHandler struct {
	usecase tagUsecase
}

func NewTagHandler(uc tagUsecase) *tagHandler {
	return &tagHandler{
		usecase: uc,
	}
}

func (h *tagHandler) List(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	userID := params.Get("userId")
	if len(userID) == 0 {
		errResponse(w, http.StatusBadRequest, "empty value 'userId'", nil)
		return
	}
	uid, err := strconv.Atoi(userID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	atags, err := h.usecase.List(uid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	res := struct {
		Tags []adapter.Tag `json:"tags"`
	}{
		Tags: atags,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}
}

func (h *tagHandler) Detail(w http.ResponseWriter, r *http.Request) {
	tagID := chi.URLParam(r, "id")
	if len(tagID) == 0 {
		errResponse(w, http.StatusBadRequest, "empty value 'tagId'", nil)
		return
	}
	tid, err := strconv.Atoi(tagID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	atag, err := h.usecase.Detail(tid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	if err := json.NewEncoder(w).Encode(atag); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}
}

func (h *tagHandler) Update(w http.ResponseWriter, r *http.Request) {
	tagID := chi.URLParam(r, "id")
	if len(tagID) == 0 {
		errResponse(w, http.StatusBadRequest, "empty value 'tagId'", nil)
		return
	}
	tid, err := strconv.Atoi(tagID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	updatedTag := adapter.Tag{
		ID: tid,
	}
	if err := json.NewDecoder(r.Body).Decode(&updatedTag); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	if err := h.usecase.Update(updatedTag); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *tagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tagID := chi.URLParam(r, "id")
	if len(tagID) == 0 {
		errResponse(w, http.StatusBadRequest, "empty value 'tagId'", nil)
		return
	}
	tid, err := strconv.Atoi(tagID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	if err := h.usecase.Delete(tid); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *tagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createTag adapter.Tag
	if err := json.NewDecoder(r.Body).Decode(&createTag); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	if err := h.usecase.Create(createTag); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
