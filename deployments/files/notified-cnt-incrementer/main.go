package p

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type afterQuery struct {
	description string
	query       string
}

var afterXdaysQueries = []*afterQuery{
	&afterQuery{
		description: "After 1 day!",
		query:       `SELECT id from memos WHERE updated_at < NOW() - interval '1 days' AND notified_cnt = 0 ORDER BY id`,
	},
	&afterQuery{
		description: "After 4 days!",
		query:       `SELECT id from memos WHERE updated_at < NOW() - interval '4 days' AND notified_cnt = 1 ORDER BY id`,
	},
	&afterQuery{
		description: "After 7 days!!",
		query:       `SELECT id from memos WHERE updated_at < NOW() - interval '7 days' AND notified_cnt = 2 ORDER BY id`,
	},
	&afterQuery{
		description: "After 11 days!!",
		query:       `SELECT id from memos WHERE updated_at < NOW() - interval '11 days' AND notified_cnt = 3 ORDER BY id`,
	},
	&afterQuery{
		description: "After 15 days!!",
		query:       `SELECT id from memos WHERE updated_at < NOW() - interval '15 days' AND notified_cnt = 4 ORDER BY id`,
	},
	&afterQuery{
		description: "After 20 days!!!",
		query:       `SELECT id from memos WHERE updated_at < NOW() - interval '20 days' AND notified_cnt = 5 ORDER BY id`,
	},
}

var (
	Conn *sql.DB
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
	}()

	log.Println("start notifiedCnt increment")
	if err := increase(); err != nil {
		return err
	}
	log.Println("success notifiedCnt increment")
	return nil
}

func increase() error {
	for _, afterQuery := range afterXdaysQueries {
		if err := execQuery(afterQuery); err != nil {
			return err
		}
	}
	return nil
}

var updateNotifiedCntQuery = "UPDATE memos SET notified_cnt = notified_cnt + 1 WHERE id IN (%s)"

func execQuery(aq *afterQuery) error {
	_, err := Conn.Exec(fmt.Sprintf(updateNotifiedCntQuery, aq.query))
	if err != nil {
		return err
	}
	return nil
}
