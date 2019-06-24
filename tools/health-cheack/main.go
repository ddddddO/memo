package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	target = "https://ddddddo.work"
	proxy  = "http://82.146.160.74:80" // ref: http://www.freeproxylists.net/ja/
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
	/*
	err = exec.Command("reboot").Run()
	if err != nil {
		log.Fatal("failed to reboot")
	}
	*/
}
