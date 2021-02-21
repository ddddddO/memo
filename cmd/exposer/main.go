package main

import (
	"flag"
	"log"
	"time"

	"github.com/ddddddO/tag-mng/exposer"
)

var (
	dsn      = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	interval = 5 * time.Minute
)

// hugo new site hogehoge で生成したhogehogeディレクトリ内でこのプログラムを実行する前提
func main() {
	flag.StringVar(&dsn, "dsn", dsn, "connection DB data source name")
	flag.DurationVar(&interval, "interval", interval, "pooling interval(ex: --interval=5m)")
	flag.Parse()

	conf := exposer.Config{
		Dsn:      dsn,
		Interval: interval,
	}

	err := exposer.Run(conf)
	if err != nil {
		log.Fatalf("failed to expose memos\n%+v", err)
	}
}
