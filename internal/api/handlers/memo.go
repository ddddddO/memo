package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

type Memo struct {
	Id      int    `json:"id"`
	Subject string `json:"subject"`
}

type Memos struct {
	MemoList []Memo `json:"memo_list"`
}

func MemoListHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		userId := params.Get("userId")
		if len(userId) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'userId'", nil)
			return
		}
		tagId := params.Get("tagId")

		var rows *sql.Rows
		var memos Memos
		var err error
		// NOTE: tagIdが設定されていない場合
		if len(tagId) == 0 {
			query := "SELECT id, subject FROM memos WHERE users_id=$1 ORDER BY id"
			rows, err = DB.Query(query, userId)
			if err != nil {
				errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
				return
			}
		} else {
			// NOTE: tagIdが設定されている場合
			query := "SELECT id, subject FROM memos WHERE users_id=$1 AND id IN (SELECT memos_id FROM memo_tag WHERE tags_id=$2) ORDER BY id"
			rows, err = DB.Query(query, userId, tagId)
			if err != nil {
				errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
				return
			}
		}

		for rows.Next() {
			var memo Memo
			if err := rows.Scan(&memo.Id, &memo.Subject); err != nil {
				errResponse(w, http.StatusInternalServerError, "failed to connect db 4", err)
				return
			}
			memos.MemoList = append(memos.MemoList, memo)
		}

		memosJson, err := json.Marshal(memos)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(memosJson))
	}
}
