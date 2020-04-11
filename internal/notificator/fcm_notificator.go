package notificator

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

type FCMNotificator struct {
	endpoint string
	token    string
	authKey  string
	dsn      string
}

type data struct {
	To           string       `json:"to"`
	Notification notification `json:"notification"`
}

type notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Icon  string `json:"icon"`
}

type afterQuery struct {
	description string
	query       string
	rsltRows    *sql.Rows
}

var afterXdaysQueries = []*afterQuery{
	&afterQuery{
		description: "After 1 day!",
		query:       `SELECT subject from memos WHERE updated_at < NOW() - interval '1 days' AND notified_cnt = 0 ORDER BY id`,
	},
	&afterQuery{
		description: "After 4 days!",
		query:       `SELECT subject from memos WHERE updated_at < NOW() - interval '4 days' AND notified_cnt = 1 ORDER BY id`,
	},
	&afterQuery{
		description: "After 7 days!!",
		query:       `SELECT subject from memos WHERE updated_at < NOW() - interval '7 days' AND notified_cnt = 2 ORDER BY id`,
	},
	&afterQuery{
		description: "After 11 days!!",
		query:       `SELECT subject from memos WHERE updated_at < NOW() - interval '11 days' AND notified_cnt = 3 ORDER BY id`,
	},
	&afterQuery{
		description: "After 15 days!!",
		query:       `SELECT subject from memos WHERE updated_at < NOW() - interval '15 days' AND notified_cnt = 4 ORDER BY id`,
	},
	&afterQuery{
		description: "After 20 days!!!",
		query:       `SELECT subject from memos WHERE updated_at < NOW() - interval '20 days' AND notified_cnt = 5 ORDER BY id`,
	},
}

func (fcmn FCMNotificator) detect() error {
	for _, afterQuery := range afterXdaysQueries {
		if err := fcmn.execQuery(afterQuery); err != nil {
			return err
		}
	}
	return nil
}

var updateNotifiedCntQuery = "UPDATE memos SET notified_cnt = notified_cnt + 1 WHERE subject IN (%s)"

func (fcmn FCMNotificator) execQuery(aq *afterQuery) error {
	conn, err := sql.Open("postgres", fcmn.dsn)
	if err != nil {
		return err
	}
	rows, err := conn.Query(aq.query)
	if err != nil {
		return err
	}
	aq.rsltRows = rows

	// NOTE: 通知対象のメモのnotified_cntを+1する
	_, err = conn.Exec(fmt.Sprintf(updateNotifiedCntQuery, aq.query))
	if err != nil {
		return err
	}
	return nil
}

func (fcmn FCMNotificator) notify() error {
	for _, afterQuery := range afterXdaysQueries {
		if err := fcmn.send(afterQuery); err != nil {
			return err
		}
	}
	return nil
}

// NOTE: 通知は１メモにつき１通ずつ送りたい。
// が、複数メモの更新日時がほとんど変わらない場合、一度にほぼ同時に連続で通知してしまう(のはいや)
// なので、1通ごとにsleepして通知する
func (fcmn FCMNotificator) send(aq *afterQuery) error {
	defer aq.rsltRows.Close()

	var (
		d   data
		cnt int
	)
	for aq.rsltRows.Next() {
		cnt++
		var (
			subject string
		)
		if err := aq.rsltRows.Scan(&subject); err != nil {
			return err
		}
		log.Println(subject)

		// FIXME: 一旦、各クエリの最後の1件だけ通知
		d = data{
			To: fcmn.token,
			Notification: notification{
				Title: aq.description,
				Body:  subject,
				Icon:  "./img/icons/android-chrome-192x192.png",
			},
		}
	}
	if cnt == 0 {
		return nil
	}

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	body := bytes.NewReader(b)
	req, err := http.NewRequest("POST", fcmn.endpoint, body)
	req.Header.Set("Authorization", fcmn.authKey)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to send notification")
	}
	// FIXME: 文字列で判定止める
	bb, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(bb), `"failure":0`) {
		return errors.New("failed to send notification")
	}

	log.Println("succeed to send notification")
	return nil
}
