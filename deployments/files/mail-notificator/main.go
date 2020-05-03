package p

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/smtp"
	"log"
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

var conn *sql.DB

func init() {
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	conn, _ = sql.Open("postgres", DBDSN)
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// XXX
func Run(ctx context.Context, m PubSubMessage) error {
	defer conn.Close()
	if err := detect(); err != nil {
		return err
	}

	if err := notify(); err != nil {
		return err
	}

	log.Println(string(m.Data))
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
		hostname = "smtp.mail.yahoo.co.jp"
		username = "lbfdeatq_0922"
		from     = "lbfdeatq_0922@yahoo.co.jp"
		to       = "lbfdeatq_0922@yahoo.co.jp"
		subject  = subject(aq.description)
		body     = body(aq.memos)
		mail     = "From: " + from + "\r\n" + "To: " + to + "\r\n" + "Subject: " + subject + "\r\n\r\n" + body
		password = os.Getenv("MAIL_PASSWORD")
	)

	tlsConn, err := tls.Dial("tcp", hostname+":465", &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         hostname,
	})
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(tlsConn, hostname)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", username, password, hostname)
	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(wc, mail)
	if err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	if err := client.Quit(); err != nil {
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
