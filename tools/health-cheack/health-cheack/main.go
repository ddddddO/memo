package main

import (
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/ddddddO/tag-mng/tools/health-cheack/health-cheack/lib"
)

const (
	target   = "https://ddddddo.work"
	proxySrc = "http://www.gatherproxy.com/ja"
	filePath = "./tmp.txt"
)

func main() {
	log.Print("start health cheack")

	pxs, err := lib.FetchProxys(proxySrc, filePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := cheackTarget(pxs, target); err != nil {
		log.Fatal(err)
	}

	log.Print("succeeded!")
}

func cheackTarget(pxs []*lib.Proxy, target string) error {
	log.Print("max attack: ", len(pxs))
	for _, px := range pxs {
		proxy := px.IP + ":" + px.Port
		log.Print("proxy:", proxy)

		cmd := exec.Command("curl", "-v", "-k", target, "-x", proxy)
		outerr, err := cmd.CombinedOutput()
		if err != nil {
			continue
		}

		dst := string(outerr)
		if strings.Contains(dst, "HTTP/1.1 200 OK") {
			return nil
		}
	}

	return errors.New("probably not accessible")
}
