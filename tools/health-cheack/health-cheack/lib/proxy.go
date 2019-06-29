package lib

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	gq "github.com/PuerkitoBio/goquery"
)

func FetchProxys(proxySrc, filePath string) ([]*Proxy, error) {
	if err := genProxysFile(proxySrc, filePath); err != nil {
		return nil, err
	}
	defer removeFile(filePath)

	pxs, err := fetchProxys(filePath)
	if err != nil {
		return nil, err
	}

	return pxs, nil
}

func genProxysFile(proxySrc, filePath string) error {
	doc, err := gq.NewDocument(proxySrc)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(doc.Text())
	if err != nil {
		return err
	}

	return nil
}

type Proxy struct {
	IP   string
	Port string
}

func fetchProxys(filePath string) ([]*Proxy, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pxs := []*Proxy{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		row := sc.Text()
		if !strings.Contains(row, "insertPrx") {
			continue
		}

		px := parseProxy(row)
		if px == nil {
			continue
		}
		pxs = append(pxs, px)
	}

	return pxs, nil
}

var r = regexp.MustCompile(`"PROXY_IP":"([0-9.]+)".*"PROXY_PORT":"([0-9A-Z]+)"`)

func parseProxy(row string) *Proxy {
	parsed := r.FindStringSubmatch(row)
	if len(parsed) != 3 {
		return nil
	}

	// 80, 8080ポートのみ
	port := ""
	switch parsed[2] {
	case "50":
		port = "80"
	case "1F90":
		port = "8080"
	default:
		return nil
	}

	return &Proxy{
		IP:   parsed[1],
		Port: port,
	}
}

func removeFile(filePath string) error {
	return nil
}
