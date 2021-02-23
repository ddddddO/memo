package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"

	_ "github.com/lib/pq"

	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

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

func MemoDetailHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memoId := chi.URLParam(r, "id")
		if len(memoId) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'memoId'", nil)
			return
		}

		params := r.URL.Query()
		userId := params.Get("userId")
		if len(userId) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'userId'", nil)
			return
		}

		const memoDetailQuery = `
	SELECT
	    m.id AS id,
	    m.subject AS subject,
		m.content AS content,
		m.is_exposed AS is_exposed,
		(SELECT jsonb_agg(t.id)
		  FROM memos m
		  JOIN memo_tag mt
		  ON m.id = mt.memos_id
		  JOIN tags t
		  ON mt.tags_id = t.id
	      WHERE m.id = $1 AND m.users_id = $2
		) AS tag_ids,
		(SELECT jsonb_agg(t.name)
		  FROM memos m
		  JOIN memo_tag mt
		  ON m.id = mt.memos_id
		  JOIN tags t
		  ON mt.tags_id = t.id
	      WHERE m.id = $1 AND m.users_id = $2
		) AS tag_names
		   FROM memos m
		   JOIN memo_tag mt
		   ON m.id = mt.memos_id
		   JOIN tags t
		   ON mt.tags_id = t.id
		WHERE m.id = $1 AND m.users_id = $2
		GROUP BY m.id
	`

		rows, err := DB.Query(memoDetailQuery, memoId, userId)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}
		rows.Next()
		var (
			memoDetail domain.Memo
			tagIds     string
			tagNames   string
		)
		// NOTE: 気持ち悪いけど、tagIds/tagNamesは別変数で取得して、sliceに変換してmemoDetailのフィールドに格納する
		err = rows.Scan(
			&memoDetail.ID, &memoDetail.Subject, &memoDetail.Content, &memoDetail.IsExposed,
			&tagIds, &tagNames,
		)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
			return
		}
		memoDetail.TagIDs = strToIntSlice(tagIds)
		memoDetail.TagNames = strToStrSlice(tagNames)

		//ref: https://qiita.com/shohei-ojs/items/311ef080cd5cff1e0e16
		var memoDetailJson bytes.Buffer
		encoder := json.NewEncoder(&memoDetailJson)
		encoder.SetEscapeHTML(false)
		encoder.Encode(memoDetail)

		w.WriteHeader(http.StatusOK)
		w.Write(memoDetailJson.Bytes())
	}
}

func strToIntSlice(s string) []int {
	if len(s) <= 2 {
		return []int{}
	}

	ss := strings.Split(s[1:len(s)-1], ",")
	var nums []int
	for _, strNum := range ss {
		strNum = strings.TrimSpace(strNum)
		num, err := strconv.Atoi(strNum)
		if err != nil {
			panic(err)
		}
		nums = append(nums, num)
	}
	return nums
}

func strToStrSlice(s string) []string {
	if len(s) <= 2 {
		return []string{}
	}

	return strings.Split(s[1:len(s)-1], ",")
}

func MemoUpdateHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----MemoUpdateHandler----")
		var updatedMemo domain.Memo
		if err := json.NewDecoder(r.Body).Decode(&updatedMemo); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to unmarshal json", err)
			return
		}

		const updateMemoQuery = `
		UPDATE memos SET subject=$1, content=$2, is_exposed=$3
		 WHERE id=$4 AND users_id=$5
		`

		result, err := DB.Exec(updateMemoQuery,
			updatedMemo.Subject, updatedMemo.Content, updatedMemo.IsExposed, updatedMemo.ID, updatedMemo.UserID,
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
			errResponse(w, http.StatusInternalServerError, "failed to update memo", nil)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func MemoCreateHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----MemoCreateHandler----")
		var createdMemo domain.Memo
		if err := json.NewDecoder(r.Body).Decode(&createdMemo); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		var createMemoQuery = `
WITH inserted AS (INSERT INTO memos(subject, content, users_id) VALUES($1, $2, $3) RETURNING id)
INSERT INTO memo_tag(memos_id, tags_id) VALUES
	((SELECT id FROM inserted), 1)
	%s;
`

		var valuesStr string
		for _, tagID := range createdMemo.TagIDs {
			valuesStr += fmt.Sprintf(",((SELECT id FROM inserted), %d)", tagID)
		}
		createMemoQuery = fmt.Sprintf(createMemoQuery, valuesStr)

		_, err := DB.Exec(createMemoQuery,
			createdMemo.Subject, createdMemo.Content, createdMemo.UserID,
		)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func MemoDeleteHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----MemoDeleteHandler----")
		var deleteMemo domain.Memo
		if err := json.NewDecoder(r.Body).Decode(&deleteMemo); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		const deleteMemoQuery = `
DELETE FROM memos WHERE users_id = $1 AND id = $2;
`

		result, err := DB.Exec(deleteMemoQuery,
			deleteMemo.UserID, deleteMemo.ID,
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
