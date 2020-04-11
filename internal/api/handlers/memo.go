package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Memo struct {
	Id      int    `json:"id"`
	Subject string `json:"subject"`
}

type Memos struct {
	MemoList []Memo `json:"memo_list"`
}

func MemoListHandler(c *gin.Context) {
	userId := c.Query("userId")
	if len(userId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "empty value 'userId'",
		})
		return
	}
	tagId := c.Query("tagId")

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

	var rows *sql.Rows
	var memos Memos
	// NOTE: tagIdが設定されていない場合
	if len(tagId) == 0 {
		query := "SELECT id, subject FROM memos WHERE users_id=$1 ORDER BY id"
		rows, err = conn.Query(query, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to connect db 2",
			})
			return
		}
	} else {
		// NOTE: tagIdが設定されている場合
		query := "SELECT id, subject FROM memos WHERE users_id=$1 AND id IN (SELECT memos_id FROM memo_tag WHERE tags_id=$2) ORDER BY id"
		rows, err = conn.Query(query, userId, tagId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to connect db 3",
			})
			return
		}
	}

	for rows.Next() {
		var memo Memo
		if err := rows.Scan(&memo.Id, &memo.Subject); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to connect db 4",
			})
			return
		}
		memos.MemoList = append(memos.MemoList, memo)
	}

	memosJson, err := json.Marshal(memos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to xxx",
		})
		return
	}

	c.JSON(http.StatusOK, string(memosJson))
}
