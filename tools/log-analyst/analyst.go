package main

import (
	"time"
)

// ログの見方
// https://knowledge.sakura.ad.jp/3424/ から
type Analyzer interface {
	collectIPAddrs() error
	collectAccessDates() error
	collectAccessPaths() error
	collectStatusCodes() error
	collectUserAgents() error
}

type H2OLogAnalyst struct {
	data         string
	analyzedDate time.Time
	ipAddrs      []string
	accessDates  []time.Time
	accessPaths  []string
	statusCodes  []string
	UserAgents   []string
}

/*
type UnicornLogAnalyst struct {
}
*/

func NewH2OLogAnalyst(data string) *H2OLogAnalyst {
	return &H2OLogAnalyst{
		data: data,
	}
}

// 解析のメイン処理
// 外枠作成
// 各々並列処理でいけるかも
func analyze(a Analyzer) error {
	err := a.collectIPAddrs()
	err = a.collectAccessDates()
	err = a.collectAccessPaths()
	err = a.collectStatusCodes()
	err = a.collectUserAgents()

	return err
}
