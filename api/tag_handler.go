package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/ddddddO/memo/api/adapter"
	"github.com/ddddddO/memo/models"
)

type tagRepository interface {
	FetchList(userID int) ([]*models.Tag, error)
	FetchListByMemoID(memoID int) ([]*models.Tag, error)
	Fetch(tagID int) (*models.Tag, error)
	Update(tag *models.Tag) error
	Delete(tagID int) error
	Create(tag *models.Tag) error
}

type tagHandler struct {
	repo tagRepository
}

func NewTagHandler(repo tagRepository) *tagHandler {
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

	atags := make([]adapter.Tag, len(tags))
	for i, tag := range tags {
		atags[i] = adapter.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		}
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
	tag, err := h.repo.Fetch(tid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	atag := adapter.Tag{
		ID:   tag.ID,
		Name: tag.Name,
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

	tag, err := h.repo.Fetch(tid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}
	tag.Name = updatedTag.Name

	if err := h.repo.Update(tag); err != nil {
		log.Println("failed to update tag", err)
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

	if err := h.repo.Delete(tid); err != nil {
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

	tag := &models.Tag{
		Name: createTag.Name,
		UsersID: sql.NullInt64{
			Int64: int64(createTag.UserID),
			Valid: true,
		},
	}
	if err := h.repo.Create(tag); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
