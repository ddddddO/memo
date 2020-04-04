package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type MemoDetail struct {
	Id       int      `json:"id"`
	Subject  string   `json:"subject"`
	Content  string   `json:"content"`
	TagIds   []int    `json:"tag_ids"`
	TagNames []string `json:"tag_names"`
}

func MemoDetailHandler(c *gin.Context) {
	memoId := c.Query("memoId")
	if len(memoId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "empty value 'memoId'",
		})
		return
	}

	userId := c.Query("userId")
	if len(userId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "empty value 'userId'",
		})
		return
	}

	// TODO: 共通化
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		log.Println("set default DSN")
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	conn, err := sql.Open("postgres", DBDSN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 1",
		})
		return
	}

	const memoDetailQuery = `
	SELECT
	    m.id AS id,
	    m.subject AS subject,
		m.content AS content,
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

	rows, err := conn.Query(memoDetailQuery, memoId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 2",
		})
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
		&memoDetail.Id, &memoDetail.Subject, &memoDetail.Content,
		&tagIds, &tagNames,
	)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 4",
		})
		return
	}
	memoDetail.TagIds = strToIntSlice(tagIds)
	memoDetail.TagNames = strToStrSlice(tagNames)

	//ref: https://qiita.com/shohei-ojs/items/311ef080cd5cff1e0e16
	var memoDetailJson bytes.Buffer
	encoder := json.NewEncoder(&memoDetailJson)
	encoder.SetEscapeHTML(false)
	encoder.Encode(memoDetail)

	c.PureJSON(http.StatusOK, memoDetailJson.String())
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
	UserId      int    `json:"user_id"`
	MemoId      int    `json:"memo_id"`
	MemoSubject string `json:"memo_subject"`
	MemoContent string `json:"memo_content"`
}

func MemoDetailUpdateHandler(c *gin.Context) {
	log.Print("----MemoDetailUpdateHandler----")
	var updatedMemo UpdatedMemo
	if err := c.BindJSON(&updatedMemo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to bind json",
		})
		return
	}

	log.Printf("%+v", updatedMemo)

	// TODO: 共通化
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		log.Println("set default DSN")
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	conn, err := sql.Open("postgres", DBDSN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 1",
		})
		return
	}

	const updateMemoQuery = `
UPDATE memos SET subject=$1, content=$2
 WHERE id=$3 AND users_id=$4
`
	result, err := conn.Exec(updateMemoQuery,
		updatedMemo.MemoSubject, updatedMemo.MemoContent, updatedMemo.MemoId, updatedMemo.UserId,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 2",
		})
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 3",
		})
		return
	}
	if n != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update memo",
		})
		return
	}
}
