package notificator

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
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

func (fcmn FCMNotificator) detect() error {
	afterXdaysQueries := []string{
		`SELECT id from memos WHERE updated_at < NOW() - interval '1 days' AND notified_cnt = 0 ORDER BY id`,
		`SELECT id from memos WHERE updated_at < NOW() - interval '4 days' AND notified_cnt = 1 ORDER BY id`,
		`SELECT id from memos WHERE updated_at < NOW() - interval '7 days' AND notified_cnt = 2 ORDER BY id`,
		`SELECT id from memos WHERE updated_at < NOW() - interval '11 days' AND notified_cnt = 3 ORDER BY id`,
		`SELECT id from memos WHERE updated_at < NOW() - interval '15 days' AND notified_cnt = 4 ORDER BY id`,
		`SELECT id from memos WHERE updated_at < NOW() - interval '20 days' AND notified_cnt = 5 ORDER BY id`,
	}
	for _, query := range afterXdaysQueries {
		if err := fcmn.execQuery(query); err != nil {
			return err
		}
	}

	return nil
}

//var updateNotifiedCntQuery = "UPDATE memos SET notified_cnt = notified_cnt + 1 WHERE id IN (%s)"

func (fcmn FCMNotificator) execQuery(query string) error {
	conn, err := sql.Open("postgres", fcmn.dsn)
	if err != nil {
		return err
	}
	rows, err := conn.Query(query)
	if err != nil {
		return err
	}

	log.Println("selected ids")
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		log.Println(id)
	}
	// NOTE: afterXdaysQueryで取得したメモを通知して、そのメモを以下なクエリでnotified_cnt++する
	// conn.Exec(fmt.Sprint(updateNotifiedCntQuery, query))
	return nil
}

// NOTE: 通知は１メモにつき１通ずつ送りたい。
// が、複数メモの更新日時がほとんど変わらない場合、一度にほぼ同時に連続で通知してしまう(のはいや)
// なので、1通ごとにsleepして通知する
func (fcmn FCMNotificator) send() error {
	d := data{
		To: fcmn.token,
		Notification: notification{
			Title: "FCM Message by go",
			Body:  "This is an FCM Message",
			Icon:  "./img/icons/android-chrome-192x192.png",
		},
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
