package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"

	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/domain"
)

type Memos struct {
	MemoList []domain.Memo `json:"memo_list"`
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
			query := "SELECT id, subject, notified_cnt FROM memos WHERE users_id=$1 ORDER BY id"
			rows, err = DB.Query(query, userId)
			if err != nil {
				errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
				return
			}
		} else {
			// NOTE: tagIdが設定されている場合
			query := "SELECT id, subject, notified_cnt FROM memos WHERE users_id=$1 AND id IN (SELECT memos_id FROM memo_tag WHERE tags_id=$2) ORDER BY id"
			rows, err = DB.Query(query, userId, tagId)
			if err != nil {
				errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
				return
			}
		}

		for rows.Next() {
			var memo domain.Memo
			if err := rows.Scan(&memo.ID, &memo.Subject, &memo.NotifiedCnt); err != nil {
				errResponse(w, http.StatusInternalServerError, "failed to connect db 4", err)
				return
			}
			setColor(&memo)
			memos.MemoList = append(memos.MemoList, memo)
		}
		// NOTE: NotifiedCntでメモを昇順にソート
		sort.SliceStable(memos.MemoList,
			func(i, j int) bool {
				return memos.MemoList[i].NotifiedCnt < memos.MemoList[j].NotifiedCnt
			},
		)

		memosJson, err := json.Marshal(memos)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(memosJson))
	}
}

func setColor(m *domain.Memo) {
	switch m.NotifiedCnt {
	case 0:
		m.RowVariant = "danger"
	case 1:
		m.RowVariant = "warning"
	case 2:
		m.RowVariant = "primary"
	case 3:
		m.RowVariant = "info"
	case 4:
		m.RowVariant = "secondary"
	case 5:
		m.RowVariant = "success"
	}
}
