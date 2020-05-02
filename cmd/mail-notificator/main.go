package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type memo struct {
	ID      int
	Subject string
}

type afterQuery struct {
	description string
	query       string
	memos       []memo
}

var afterXdaysQueries = []*afterQuery{
	&afterQuery{
		description: "After 1 day!",
		query:       `SELECT id, subject from memos WHERE updated_at < NOW() - interval '1 days' AND notified_cnt = 0 ORDER BY id`,
	},
	&afterQuery{
		description: "After 4 days!",
		query:       `SELECT id, subject from memos WHERE updated_at < NOW() - interval '4 days' AND notified_cnt = 1 ORDER BY id`,
	},
	&afterQuery{
		description: "After 7 days!!",
		query:       `SELECT id, subject from memos WHERE updated_at < NOW() - interval '7 days' AND notified_cnt = 2 ORDER BY id`,
	},
	&afterQuery{
		description: "After 11 days!!",
		query:       `SELECT id, subject from memos WHERE updated_at < NOW() - interval '11 days' AND notified_cnt = 3 ORDER BY id`,
	},
	&afterQuery{
		description: "After 15 days!!",
		query:       `SELECT id, subject from memos WHERE updated_at < NOW() - interval '15 days' AND notified_cnt = 4 ORDER BY id`,
	},
	&afterQuery{
		description: "After 20 days!!!",
		query:       `SELECT id, subject from memos WHERE updated_at < NOW() - interval '20 days' AND notified_cnt = 5 ORDER BY id`,
	},
}

// exec: MAIL_PASSWORD=XXXXX DBDSN=YYYYY go run main.go
func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	if err := detect(); err != nil {
		return err
	}

	if err := notify(); err != nil {
		return err
	}
	return nil
}

func detect() error {
	for _, afterQuery := range afterXdaysQueries {
		if err := execQuery(afterQuery); err != nil {
			return err
		}
	}
	return nil
}

func notify() error {
	for _, afterQuery := range afterXdaysQueries {
		// 検知したメモが0件の場合、メールしない
		if len(afterQuery.memos) == 0 {
			continue
		}
		if err := send(afterQuery); err != nil {
			return err
		}
	}
	return nil
}

var updateNotifiedCntQuery = "UPDATE memos SET notified_cnt = notified_cnt + 1 WHERE subject IN (%s)"

func execQuery(aq *afterQuery) error {
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	conn, err := sql.Open("postgres", DBDSN)
	if err != nil {
		return err
	}
	rows, err := conn.Query(aq.query)
	if err != nil {
		return err
	}

	for rows.Next() {
		m := memo{}
		if err := rows.Scan(&m.ID, &m.Subject); err != nil {
			return err
		}
		aq.memos = append(aq.memos, m)
	}

	// NOTE: 通知対象のメモのnotified_cntを+1する
	_, err = conn.Exec(fmt.Sprintf(updateNotifiedCntQuery, strings.Replace(aq.query, " id,", "", 1)))
	if err != nil {
		return err
	}
	return nil
}

func send(aq *afterQuery) error {
	var (
		hostname = "smtp.gmail.com"
		from     = "lbfdeatq@gmail.com"
		to       = "lbfdeatq@gmail.com"
		subject  = subject(aq.description)
		body     = body(aq.memos)
		mail     = []byte(
			"To: " + to + "\r\n" +
				"Subject: " + subject + "\r\n\r\n" +
				body,
		)
		recipients = []string{to}
		password   = os.Getenv("MAIL_PASSWORD")
	)

	auth := smtp.PlainAuth("", from, password, hostname)
	err := smtp.SendMail(hostname+":587", auth, from, recipients, mail)
	if err != nil {
		return err
	}

	return nil
}

func subject(description string) string {
	return description + " by Tag-Mng"
}

func body(memos []memo) string {
	msg := "After Login, Confirm Memo List!\r\n" + "https://XXXXXX/" + "\r\n\r\n"
	for _, memo := range memos {
		msg = msg +
			memo.Subject + "\r\n" +
			"https://XXXXXX/memodetail/" + strconv.Itoa(memo.ID) + "\r\n\r\n"
	}
	return msg
}
