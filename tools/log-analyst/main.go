package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	logPath string
)

func main() {
	//flag.StringVar(&logPath, "p", "/var/log/h2o/access_log", "アクセスログのパスを指定する") ←実装はこちら
	flag.StringVar(&logPath, "p", "../data/access_log", "アクセスログのパスを指定する")
	flag.Parse()

	fmt.Println(logPath)

	data, err := readFile(logPath)
	if err != nil {
		panic(err)
	}
	//fmt.Println(data)

	h2oLogAnalyst := NewH2OLogAnalyst(data)
	err = analyze(h2oLogAnalyst) // メソッド未実装のため怒られている
	if err != nil {
		panic(err)
	}
}

// ファイル内容をbufに格納してすぐ閉じる
func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	finf, err := f.Stat()
	if err != nil {
		return "", err
	}

	buf := make([]byte, finf.Size())
	_, err = f.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
