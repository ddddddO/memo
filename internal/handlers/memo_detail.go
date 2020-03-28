package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type MemoDetail struct {
	Id       int    `json:"id"`
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	TagIds   string `json:"tag_ids"`
	TagNames string `json:"tag_names"`
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
	SELECT DISTINCT
	    m.id AS id,
	    m.subject AS subject,
		m.content AS content,
		ARRAY(SELECT t.id
		  FROM memos m 
		  JOIN memo_tag mt 
		  ON m.id = mt.memos_id
		  JOIN tags t
		  ON mt.tags_id = t.id
	      WHERE m.id = $1 AND m.users_id = $2
		)  AS tag_ids,
		ARRAY(SELECT t.name
		  FROM memos m 
		  JOIN memo_tag mt 
		  ON m.id = mt.memos_id
		  JOIN tags t
		  ON mt.tags_id = t.id
	      WHERE m.id = $1 AND m.users_id = $2
		)  AS tag_names
		   FROM memos m 
		   JOIN memo_tag mt 
		   ON m.id = mt.memos_id
		   JOIN tags t
		   ON mt.tags_id = t.id
		WHERE m.id = $1 AND m.users_id = $2
	`

	rows, err := conn.Query(memoDetailQuery, memoId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 2",
		})
		return
	}
	rows.Next()
	var memoDetail MemoDetail
	err = rows.Scan(
		&memoDetail.Id, &memoDetail.Subject, &memoDetail.Content,
		&memoDetail.TagIds, &memoDetail.TagNames,
	)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 4",
		})
		return
	}

	//ref: https://qiita.com/shohei-ojs/items/311ef080cd5cff1e0e16
	var memoDetailJson bytes.Buffer
	encoder := json.NewEncoder(&memoDetailJson)
	encoder.SetEscapeHTML(false)
	encoder.Encode(memoDetail)

	c.PureJSON(http.StatusOK, memoDetailJson.String())
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
