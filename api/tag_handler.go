package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/domain"
	"github.com/ddddddO/tag-mng/repository"
)

func TagListHandler(repo repository.TagRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		tags, err := repo.FetchList(uid)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		res := struct {
			Tags []domain.Tag `json:"tags"`
		}{
			Tags: tags,
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}
	}
}

func TagDetailHandler(repo repository.TagRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		tag, err := repo.Fetch(tid)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		if err := json.NewEncoder(w).Encode(tag); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}
	}
}

func TagUpdateHandler(repo repository.TagRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		updatedTag := domain.Tag{
			ID: tid,
		}
		if err := json.NewDecoder(r.Body).Decode(&updatedTag); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		if err := repo.Update(updatedTag); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func TagDeleteHandler(repo repository.TagRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		deleteTag := domain.Tag{
			ID: tid,
		}
		if err := repo.Delete(deleteTag); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 1", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func TagCreateHandler(repo repository.TagRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createTag domain.Tag
		if err := json.NewDecoder(r.Body).Decode(&createTag); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		if err := repo.Create(createTag); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
