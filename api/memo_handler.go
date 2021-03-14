package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi"

	"github.com/ddddddO/tag-mng/domain"
)

type Memos struct {
	MemoList []domain.Memo `json:"memo_list"`
}

func MemoListHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		userID := params.Get("userId")
		if len(userID) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'userId'", nil)
			return
		}
		tagID := params.Get("tagId")

		var rows *sql.Rows
		var memos Memos
		var err error
		// NOTE: tagIdが設定されていない場合
		if len(tagID) == 0 {
			query := "SELECT id, subject, notified_cnt FROM memos WHERE users_id=$1 ORDER BY id"
			rows, err = db.Query(query, userID)
			if err != nil {
				errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
				return
			}
		} else {
			// NOTE: tagIdが設定されている場合
			query := "SELECT id, subject, notified_cnt FROM memos WHERE users_id=$1 AND id IN (SELECT memos_id FROM memo_tag WHERE tags_id=$2) ORDER BY id"
			rows, err = db.Query(query, userID, tagID)
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

func MemoDetailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memoID := chi.URLParam(r, "id")
		if len(memoID) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'memoId'", nil)
			return
		}

		params := r.URL.Query()
		userID := params.Get("userId")
		if len(userID) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'userId'", nil)
			return
		}

		// TODO: 見直す
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

		rows, err := db.Query(memoDetailQuery, memoID, userID)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}
		rows.Next()
		var (
			memoDetail domain.Memo
			tagIDs     string
			tagNames   string
		)
		// NOTE: 気持ち悪いけど、tagIds/tagNamesは別変数で取得して、sliceに変換してmemoDetailのフィールドに格納する
		err = rows.Scan(
			&memoDetail.ID, &memoDetail.Subject, &memoDetail.Content, &memoDetail.IsExposed,
			&tagIDs, &tagNames,
		)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
			return
		}

		memoDetail.Tags = toTags(tagIDs, tagNames)
		//ref: https://qiita.com/shohei-ojs/items/311ef080cd5cff1e0e16
		var memoDetailJson bytes.Buffer
		encoder := json.NewEncoder(&memoDetailJson)
		encoder.SetEscapeHTML(false)
		encoder.Encode(memoDetail)

		w.WriteHeader(http.StatusOK)
		w.Write(memoDetailJson.Bytes())
	}
}

func toTags(ids, names string) []domain.Tag {
	convertedIDs := strToIntSlice(ids)
	convertedNames := strToStrSlice(names)

	var tags []domain.Tag
	for i := range convertedIDs {
		tag := domain.Tag{}
		tag.ID = convertedIDs[i]
		tag.Name = convertedNames[i]

		tags = append(tags, tag)
	}

	return tags
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

func MemoUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----MemoUpdateHandler----")

		memoID := chi.URLParam(r, "id")
		if len(memoID) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'memoId'", nil)
			return
		}

		var updatedMemo domain.Memo
		if err := json.NewDecoder(r.Body).Decode(&updatedMemo); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to unmarshal json", err)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 0", err)
			return
		}
		defer tx.Rollback()

		const updateMemoQuery = `
		UPDATE memos SET subject=$1, content=$2, is_exposed=$3
		 WHERE id=$4 AND users_id=$5
		`

		result, err := tx.Exec(updateMemoQuery,
			updatedMemo.Subject, updatedMemo.Content, updatedMemo.IsExposed, memoID, updatedMemo.UserID,
		)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 1", err)
			return
		}
		n, err := result.RowsAffected()
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}
		if n != 1 {
			errResponse(w, http.StatusInternalServerError, "failed to update memo", nil)
			return
		}

		const deleteMemoTagQuery = `
		DELETE FROM memo_tag WHERE memos_id=$1 AND tags_id <> 1
		`

		_, err = tx.Exec(deleteMemoTagQuery, memoID)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
			return
		}

		var insertMemoTagQuery = `
		INSERT INTO memo_tag(memos_id, tags_id) VALUES
		%s
		`

		var valuesStr string
		for _, tag := range updatedMemo.Tags {
			valuesStr += fmt.Sprintf("(%v, %d),", memoID, tag.ID)
		}
		insertMemoTagQuery = fmt.Sprintf(insertMemoTagQuery, valuesStr[:len(valuesStr)-1])

		_, err = tx.Exec(insertMemoTagQuery)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 4", err)
			return
		}

		tx.Commit()

		w.WriteHeader(http.StatusOK)
	}
}

func MemoCreateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----MemoCreateHandler----")
		var createdMemo domain.Memo
		if err := json.NewDecoder(r.Body).Decode(&createdMemo); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		var createMemoQuery = `
WITH inserted AS (INSERT INTO memos(subject, content, users_id, is_exposed) VALUES($1, $2, $3, $4) RETURNING id)
INSERT INTO memo_tag(memos_id, tags_id) VALUES
	((SELECT id FROM inserted), 1)
	%s;
`

		var valuesStr string
		for _, tag := range createdMemo.Tags {
			valuesStr += fmt.Sprintf(",((SELECT id FROM inserted), %d)", tag.ID)
		}
		createMemoQuery = fmt.Sprintf(createMemoQuery, valuesStr)

		_, err := db.Exec(createMemoQuery,
			createdMemo.Subject, createdMemo.Content, createdMemo.UserID, createdMemo.IsExposed,
		)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func MemoDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----MemoDeleteHandler----")

		memoID := chi.URLParam(r, "id")
		if len(memoID) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'memoId'", nil)
			return
		}

		var deleteMemo domain.Memo
		if err := json.NewDecoder(r.Body).Decode(&deleteMemo); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		const deleteMemoQuery = `
DELETE FROM memos WHERE users_id = $1 AND id = $2;
`

		tx, err := db.Begin()
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 1", err)
			return
		}
		defer tx.Rollback()

		result, err := tx.Exec(deleteMemoQuery,
			deleteMemo.UserID, memoID,
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

		const deleteMemoTagQuery = `
DELETE FROM memo_tag WHERE memos_id = $1
`

		_, err = tx.Exec(deleteMemoTagQuery, memoID)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 4", err)
			return
		}

		tx.Commit()

		w.WriteHeader(http.StatusOK)
	}
}
