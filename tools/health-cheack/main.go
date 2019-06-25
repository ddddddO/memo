package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	target = "https://ddddddo.work"
	proxy  = "http://173.82.173.110:8080" // ref: http://www.freeproxylists.net/ja/
)

func main() {
	log.Print("start health cheack")

	cmd := exec.Command("curl", "-v", "-k", target, "-x", proxy)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(string(out), "HTTP/1.1 200 OK") {
		log.Print("200 OK")
		os.Exit(0)
	}
	
	log.Fatal("!200")
}
