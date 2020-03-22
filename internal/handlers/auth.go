package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func AuthHandler(c *gin.Context) {
	// TODO: name -> email へ変更したい
	name, ok := c.GetPostForm("name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "empty key 'name'",
		})
		return
	}
	passwd, ok := c.GetPostForm("passwd")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "empty key 'passwd'",
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

	const query = "SELECT id FROM users WHERE name=$1 AND passwd=$2"
	rows, err := conn.Query(query, name, genSecuredPasswd(passwd, name))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 2",
		})
		return
	}

	// TODO: ユーザー登録をしてもらう or 正しいname/passwdを指定してもらう
	if !rows.Next() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "faild to authenticate",
		})
		return
	}
	var userId int
	if err := rows.Scan(&userId); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 3",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userId,
	})
}

func genSecuredPasswd(name, passwd string) string {
	secStrPass := name + passwd
	var secPass [32]byte
	for i := 0; i < 100000; i++ {
		secPass = sha256.Sum256([]byte(secStrPass))
		secStrPass = hex.EncodeToString(secPass[:])
	}
	return strings.ToLower(hex.EncodeToString(secPass[:]))
}
