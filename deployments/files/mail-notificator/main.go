package p

import (
	"context"
	"crypto/tls"
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

var (
	Conn       *sql.DB
	TlsConn    *tls.Conn
	MailClient *smtp.Client
)

var (
	hostname = "smtp.mail.yahoo.co.jp"
	username = "lbfdeatq_0922"
	from     = "lbfdeatq_0922@yahoo.co.jp"
	to       = "lbfdeatq_0922@yahoo.co.jp"
	password = os.Getenv("MAIL_PASSWORD")
)

func init() {
	// DB初期化
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}
	var err error
	Conn, err = sql.Open("postgres", DBDSN)
	if err != nil {
		log.Fatal(err)
	}

	// mail client初期化
	TlsConn, err = tls.Dial("tcp", hostname+":465", &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         hostname,
	})
	if err != nil {
		log.Fatal(err)
	}
	MailClient, err = smtp.NewClient(TlsConn, hostname)
	if err != nil {
		log.Fatal(err)
	}
	auth := smtp.PlainAuth("", username, password, hostname)
	if err := MailClient.Auth(auth); err != nil {
		log.Fatal(err)
	}
}

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// XXX
func Run(ctx context.Context, m PubSubMessage) error {
	defer func() {
		if Conn != nil {
			Conn.Close()
		}
		if TlsConn != nil {
			TlsConn.Close()
		}
		if MailClient != nil {
			MailClient.Quit()
			MailClient.Close()
		}
	}()

	log.Println("start mail notice")
	if err := detect(); err != nil {
		return err
	}

	if err := notify(); err != nil {
		return err
	}
	log.Println("success mail notice")
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
	rows, err := Conn.Query(aq.query)
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
	// FIXME: メール出来たメモのみ更新するようにすべき
	_, err = Conn.Exec(fmt.Sprintf(updateNotifiedCntQuery, strings.Replace(aq.query, " id,", "", 1)))
	if err != nil {
		return err
	}
	return nil
}

func send(aq *afterQuery) error {
	if MailClient == nil {
		return nil
	}

	var (
		subject = subject(aq.description)
		body    = body(aq.memos)
		mail    = "From: " + from + "\r\n" + "To: " + to + "\r\n" + "Subject: " + subject + "\r\n\r\n" + body
	)

	if err := MailClient.Mail(from); err != nil {
		return err
	}
	if err := MailClient.Rcpt(to); err != nil {
		return err
	}
	wc, err := MailClient.Data()
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

	return nil
}

func subject(description string) string {
	return description + " by Tag-Mng"
}

func body(memos []memo) string {
	msg := "After Login, Confirm Memo List!\r\n" + "https://app-dot-tag-mng-243823.appspot.com/" + "\r\n\r\n"
	for _, memo := range memos {
		msg = msg +
			memo.Subject + "\r\n" +
			"https://app-dot-tag-mng-243823.appspot.com/memodetail/" + strconv.Itoa(memo.ID) + "\r\n\r\n"
	}
	return msg
}
