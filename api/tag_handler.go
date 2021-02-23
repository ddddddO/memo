package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/domain"
)

type Tags struct {
	TagList []domain.Tag `json:"tag_list"`
}

func TagListHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		userID := params.Get("userId")
		if len(userID) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'userId'", nil)
			return
		}

		var rows *sql.Rows
		var tags Tags
		var err error
		query := "SELECT id, name FROM tags WHERE users_id = $1 ORDER BY id"
		rows, err = DB.Query(query, userID)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}

		for rows.Next() {
			var tag domain.Tag
			if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
				errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
				return
			}
			tags.TagList = append(tags.TagList, tag)
		}

		tagsJson, err := json.Marshal(tags)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tagsJson))
	}
}

func TagDetailHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagId := chi.URLParam(r, "id")
		if len(tagId) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'tagId'", nil)
			return
		}

		var rows *sql.Rows
		var err error
		query := "SELECT id, name FROM tags WHERE id = $1"
		rows, err = DB.Query(query, tagId)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}

		rows.Next()
		var tag domain.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
			return
		}

		tagJson, err := json.Marshal(tag)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tagJson))
	}
}

func TagUpdateHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----TagUpdateHandler----")
		var updatedTag domain.Tag
		if err := json.NewDecoder(r.Body).Decode(&updatedTag); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		const updateTagQuery = `
UPDATE tags SET name = $1 WHERE id = $2
`
		result, err := DB.Exec(updateTagQuery,
			updatedTag.Name, updatedTag.ID,
		)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}
		n, err := result.RowsAffected()
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
			return
		}
		if n != 1 {
			errResponse(w, http.StatusInternalServerError, "failed", nil)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func TagDeleteHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----TagDeleteHandler----")
		var deleteTag domain.Tag
		if err := json.NewDecoder(r.Body).Decode(&deleteTag); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		const deleteTagQuery = `
DELETE FROM tags WHERE id = $1
`
		result, err := DB.Exec(deleteTagQuery,
			deleteTag.ID,
		)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}
		n, err := result.RowsAffected()
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
			return
		}
		if n != 1 {
			errResponse(w, http.StatusInternalServerError, "failed", nil)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func TagCreateHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----TagCreateHandler----")
		var createTag domain.Tag
		if err := json.NewDecoder(r.Body).Decode(&createTag); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		const createTagQuery = `
INSERT INTO tags(name, users_id) VALUES($1, $2) RETURNING id
`
		result, err := DB.Exec(createTagQuery,
			createTag.Name, createTag.UserID,
		)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}
		n, err := result.RowsAffected()
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
			return
		}
		if n != 1 {
			errResponse(w, http.StatusInternalServerError, "failed", nil)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
