package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type Pinger struct {
	timeout  time.Duration
	interval time.Duration
}

const defaultTimeout = 2000 * time.Millisecond
const defaultInterval = 0 * time.Millisecond

func NewPinger() *Pinger {
	return &Pinger{
		defaultTimeout,
		defaultInterval,
	}
}
func (p *Pinger) Ping(target string) (success bool, latency int) {
	begin := time.Now()
	conn, err := net.DialTimeout("tcp", target, p.timeout)
	if err != nil {
		// fmt.Println("connectoin failed.")
		return false, 0
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	end := time.Now()
	latency = int(end.Sub(begin)) / int(time.Millisecond)
	return true, latency
}
func (p *Pinger) DoPing(target string, count int) {
	host := strings.SplitN(target, ":", 2)[0]
	//prepare for IP address
	_, err := net.LookupIP(host)
	if err != nil {
		fmt.Printf("%s - latency: name resolve failed\n", target)
		return
	}
	for i := 0; i < count; i++ {
		ok, latency := p.Ping(target)
		if ok {
			fmt.Printf("%s - latency: %dms\n", target, latency)
		} else {
			fmt.Printf("%s - latency: timeout\n", target)
		}
		time.Sleep(p.interval)
	}
}

var addrExpr = regexp.MustCompile(`"add":"([\w\.-]+)"`)
var portExpr = regexp.MustCompile(`"port":(\d+)`)

func (p *Pinger) Batch(vmessLinksFile string) {
	var targets []string
	f, err := os.Open(vmessLinksFile)
	if err != nil {
		fmt.Println("Failed to open file:", vmessLinksFile)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "vmess://") {
			data, _ := base64.StdEncoding.DecodeString(line[8:])
			jsonstr := string(data)
			host := addrExpr.FindStringSubmatch(jsonstr)[1]
			port := portExpr.FindStringSubmatch(jsonstr)[1]
			targets = append(targets, fmt.Sprintf("%s:%s", host, port))
		}
	}
	fmt.Printf("found %d targets\n\n", len(targets))
	for _, target := range targets {
		p.DoPing(target, 1)
	}

}

func main() {

	pinger := NewPinger()
	// pinger.DoPing("baidu.com:80", 3)
	pinger.Batch("vmesslinks.txt")
}
