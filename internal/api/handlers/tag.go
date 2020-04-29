package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Tags struct {
	TagList []Tag `json:"tag_list"`
}

func TagListHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	userId := params.Get("userId")
	if len(userId) == 0 {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"message": "empty value 'userId'",
		// })
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
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 1",
		// })
		return
	}

	var rows *sql.Rows
	var tags Tags
	query := "SELECT id, name FROM tags WHERE users_id = $1 ORDER BY id"
	rows, err = conn.Query(query, userId)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 2",
		// })
		return
	}

	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			log.Println(err)
			// c.JSON(http.StatusInternalServerError, gin.H{
			// 	"message": "failed to connect db 4",
			// })
			return
		}
		tags.TagList = append(tags.TagList, tag)
	}

	tagsJson, err := json.Marshal(tags)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to xxx",
		// })
		return
	}

	//c.JSON(http.StatusOK, string(tagsJson))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tagsJson))
}

func TagDetailHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	tagId := params.Get("tagId")
	if len(tagId) == 0 {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"message": "empty value 'userId'",
		// })
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
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 1",
		// })
		return
	}

	var rows *sql.Rows
	query := "SELECT id, name FROM tags WHERE id = $1"
	rows, err = conn.Query(query, tagId)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 2",
		// })
		return
	}

	rows.Next()
	var tag Tag
	if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
		log.Println(err)
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 4",
		// })
		return
	}

	tagJson, err := json.Marshal(tag)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to xxx",
		// })
		return
	}

	//c.JSON(http.StatusOK, string(tagJson))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tagJson))
}

type UpdatedTag struct {
	Id   int    `json:"tag_id"`
	Name string `json:"tag_name"`
}

func TagDetailUpdateHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("----TagDetailUpdateHandler----")
	var updatedTag UpdatedTag
	buff := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buff)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(buff, updatedTag); err != nil {
		panic(err)
		return
	}
	// if err := c.BindJSON(&updatedTag); err != nil {
	// 	// c.JSON(http.StatusInternalServerError, gin.H{
	// 	// 	"message": "failed to bind json",
	// 	// })
	// 	return
	// }

	log.Printf("%+v", updatedTag)

	// TODO: 共通化
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		log.Println("set default DSN")
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	conn, err := sql.Open("postgres", DBDSN)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 1",
		// })
		return
	}

	const updateTagQuery = `
UPDATE tags SET name = $1 WHERE id = $2
`
	result, err := conn.Exec(updateTagQuery,
		updatedTag.Name, updatedTag.Id,
	)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 2",
		// })
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 3",
		// })
		return
	}
	if n != 1 {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to update memo",
		// })
		return
	}
}

type DeleteTag struct {
	Id int `json:"tag_id"`
}

func TagDetailDeleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("----TagDetailDeleteHandler----")
	var deleteTag DeleteTag
	buff := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buff)
	if err != nil {
		panic(err)
		return
	}
	if err := json.Unmarshal(buff, &deleteTag); err != nil {
		panic(err)
		return
	}
	// if err := c.BindJSON(&deleteTag); err != nil {
	// 	// c.JSON(http.StatusInternalServerError, gin.H{
	// 	// 	"message": "failed to bind json",
	// 	// })
	// 	return
	// }

	log.Printf("%+v", deleteTag)

	// TODO: 共通化
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		log.Println("set default DSN")
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	conn, err := sql.Open("postgres", DBDSN)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 1",
		// })
		return
	}

	const deleteTagQuery = `
DELETE FROM tags WHERE id = $1
`
	result, err := conn.Exec(deleteTagQuery,
		deleteTag.Id,
	)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 2",
		// })
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 3",
		// })
		return
	}
	if n != 1 {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to update memo",
		// })
		return
	}
}

type CreateTag struct {
	Id     int    `json:"tag_id"`
	Name   string `json:"tag_name"`
	UserId int    `json:"user_id"`
}

func TagDetailCreateHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("----TagDetailCreateHandler----")
	var createTag CreateTag
	buff := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buff)
	if err != nil {
		panic(err)
		return
	}
	if err := json.Unmarshal(buff, &createTag); err != nil {
		panic(err)
	}
	// if err := c.BindJSON(&createTag); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "failed to bind json",
	// 	})
	// 	return
	// }

	log.Printf("%+v", createTag)

	// TODO: 共通化
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		log.Println("set default DSN")
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	conn, err := sql.Open("postgres", DBDSN)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 1",
		// })
		return
	}

	const createTagQuery = `
INSERT INTO tags(name, users_id) VALUES($1, $2) RETURNING id
`
	result, err := conn.Exec(createTagQuery,
		createTag.Name, createTag.UserId,
	)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 2",
		// })
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 3",
		// })
		return
	}
	if n != 1 {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to update memo",
		// })
		return
	}
}
