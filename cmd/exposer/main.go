package main

import (
	"flag"
	"log"

	"github.com/ddddddO/tag-mng/exposer"
)

var (
	dsn = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
)

// hugo new site hogehoge で生成したhogehogeディレクトリ内でこのプログラムを実行する前提
func main() {
	flag.StringVar(&dsn, "dsn", dsn, "connection DB data source name")
	flag.Parse()

	err := exposer.Run(dsn)
	if err != nil {
		log.Fatalf("failed to expose memo\n%+v", err)
	}
}
