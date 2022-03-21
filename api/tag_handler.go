package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/ddddddO/memo/adapter"
	"github.com/ddddddO/memo/repository"
)

type tagHandler struct {
	repo repository.TagRepository
}

func NewTagHandler(repo repository.TagRepository) *tagHandler {
	return &tagHandler{
		repo: repo,
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

	tags, err := h.repo.FetchList(uid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	res := struct {
		Tags []adapter.Tag `json:"tags"`
	}{
		Tags: tags,
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
	tag, err := h.repo.Fetch(tid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	if err := json.NewEncoder(w).Encode(tag); err != nil {
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

	if err := h.repo.Update(updatedTag); err != nil {
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
		errResponse(w, http.StatusInternalServerError, "failed to connect db 1", err)
		return
	}

	deleteTag := adapter.Tag{
		ID: tid,
	}
	if err := h.repo.Delete(deleteTag); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed to connect db 1", err)
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

	if err := h.repo.Create(createTag); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
