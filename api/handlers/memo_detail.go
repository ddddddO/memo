package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type MemoDetail struct {
	Id        int      `json:"id"`
	Subject   string   `json:"subject"`
	Content   string   `json:"content"`
	IsExposed bool     `json:"is_exposed"`
	TagIds    []int    `json:"tag_ids"`
	TagNames  []string `json:"tag_names"`
}

func MemoDetailHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		memoId := params.Get("memoId")
		if len(memoId) == 0 {
			errResponse(w, http.StatusBadRequest, "empty value 'memoId'", nil)
			return
		}

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
			memoDetail MemoDetail
			tagIds     string
			tagNames   string
		)
		// NOTE: 気持ち悪いけど、tagIds/tagNamesは別変数で取得して、sliceに変換してmemoDetailのフィールドに格納する
		err = rows.Scan(
			&memoDetail.Id, &memoDetail.Subject, &memoDetail.Content, &memoDetail.IsExposed,
			&tagIds, &tagNames,
		)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 3", err)
			return
		}
		memoDetail.TagIds = strToIntSlice(tagIds)
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

type UpdatedMemo struct {
	UserId        int    `json:"user_id"`
	MemoId        int    `json:"memo_id"`
	MemoSubject   string `json:"memo_subject"`
	MemoContent   string `json:"memo_content"`
	MemoIsExposed bool   `json:"memo_is_exposed"`
}

func MemoDetailUpdateHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----MemoDetailUpdateHandler----")
		var updatedMemo UpdatedMemo
		buff := make([]byte, r.ContentLength)
		_, err := r.Body.Read(buff)
		if err != nil && err != io.EOF {
			panic(err)
			return
		}
		if err := json.Unmarshal(buff, &updatedMemo); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to unmarshal json", err)
			return
		}

		const updateMemoQuery = `
		UPDATE memos SET subject=$1, content=$2, is_exposed=$3
		 WHERE id=$4 AND users_id=$5
		`

		result, err := DB.Exec(updateMemoQuery,
			updatedMemo.MemoSubject, updatedMemo.MemoContent, updatedMemo.MemoIsExposed, updatedMemo.MemoId, updatedMemo.UserId,
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

type CreatedMemo struct {
	UserId      int    `json:"user_id"`
	TagIds      []int  `json:"tag_ids"`
	MemoSubject string `json:"memo_subject"`
	MemoContent string `json:"memo_content"`
}

func MemoDetailCreateHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----MemoDetailCreateHandler----")
		var createdMemo CreatedMemo
		buff := make([]byte, r.ContentLength)
		_, err := r.Body.Read(buff)
		if err != nil && err != io.EOF {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}
		if err := json.Unmarshal(buff, &createdMemo); err != nil {
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
		for _, tagId := range createdMemo.TagIds {
			valuesStr += fmt.Sprintf(",((SELECT id FROM inserted), %d)", tagId)
		}
		createMemoQuery = fmt.Sprintf(createMemoQuery, valuesStr)

		_, err = DB.Exec(createMemoQuery,
			createdMemo.MemoSubject, createdMemo.MemoContent, createdMemo.UserId,
		)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type DeleteMemo struct {
	UserId int `json:"user_id"`
	MemoId int `json:"memo_id"`
}

func MemoDetailDeleteHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("----MemoDetailDeleteHandler----")
		var deleteMemo DeleteMemo
		buff := make([]byte, r.ContentLength)
		_, err := r.Body.Read(buff)
		if err != nil && err != io.EOF {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}
		if err := json.Unmarshal(buff, &deleteMemo); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		const deleteMemoQuery = `
DELETE FROM memos WHERE users_id = $1 AND id = $2;
`

		result, err := DB.Exec(deleteMemoQuery,
			deleteMemo.UserId, deleteMemo.MemoId,
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
