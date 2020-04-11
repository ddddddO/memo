package notificator

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type FCMNotificator struct {
	endpoint string
	token    string
	authKey  string
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
	log.Println("DETECT")
	return nil
}

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
