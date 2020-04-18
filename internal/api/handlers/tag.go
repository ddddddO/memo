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

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Tags struct {
	TagList []Tag `json:"tag_list"`
}

func TagListHandler(c *gin.Context) {
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

	var rows *sql.Rows
	var tags Tags
	query := "SELECT id, name FROM tags WHERE users_id = $1 ORDER BY id"
	rows, err = conn.Query(query, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 2",
		})
		return
	}

	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to connect db 4",
			})
			return
		}
		tags.TagList = append(tags.TagList, tag)
	}

	tagsJson, err := json.Marshal(tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to xxx",
		})
		return
	}

	c.JSON(http.StatusOK, string(tagsJson))
}

func TagDetailHandler(c *gin.Context) {
	tagId := c.Query("tagId")
	if len(tagId) == 0 {
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

	var rows *sql.Rows
	query := "SELECT id, name FROM tags WHERE id = $1"
	rows, err = conn.Query(query, tagId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 2",
		})
		return
	}

	rows.Next()
	var tag Tag
	if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 4",
		})
		return
	}

	tagJson, err := json.Marshal(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to xxx",
		})
		return
	}

	c.JSON(http.StatusOK, string(tagJson))
}

type UpdatedTag struct {
	Id   int    `json:"tag_id"`
	Name string `json:"tag_name"`
}

func TagDetailUpdateHandler(c *gin.Context) {
	log.Print("----TagDetailUpdateHandler----")
	var updatedTag UpdatedTag
	if err := c.BindJSON(&updatedTag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to bind json",
		})
		return
	}

	log.Printf("%+v", updatedTag)

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

	const updateTagQuery = `
UPDATE tags SET name = $1 WHERE id = $2
`
	result, err := conn.Exec(updateTagQuery,
		updatedTag.Name, updatedTag.Id,
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

type DeleteTag struct {
	Id int `json:"tag_id"`
}

func TagDetailDeleteHandler(c *gin.Context) {
	log.Print("----TagDetailDeleteHandler----")
	var deleteTag DeleteTag
	if err := c.BindJSON(&deleteTag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to bind json",
		})
		return
	}

	log.Printf("%+v", deleteTag)

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

	const deleteTagQuery = `
DELETE FROM tags WHERE id = $1
`
	result, err := conn.Exec(deleteTagQuery,
		deleteTag.Id,
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
