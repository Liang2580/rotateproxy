package main

import (
	"flag"

	"github.com/akkuman/rotateproxy"
)

var (
	baseCfg   rotateproxy.BaseConfig
	email     string
	token     string
	rule      string
	pageCount int
	timeout int
)

func init() {
	flag.StringVar(&email, "email", "", "email address")
	flag.StringVar(&token, "token", "", "token")
	flag.StringVar(&baseCfg.ListenAddr, "l", ":9999", "listen address")
	flag.IntVar(&timeout, "time", 8000, "Not used for more than milliseconds")
	// && country="CN"
	flag.StringVar(&rule, "rule", `protocol=="socks5" && "Version:5 Method:No Authentication(0x00)" && after="2021-12-01"`, "search rule")
	flag.IntVar(&baseCfg.IPRegionFlag, "region", 0, "0: all 1: cannot bypass gfw 2: bypass gfw")
	flag.IntVar(&pageCount, "page", 5, "the page count you want to crawl")
	flag.Parse()
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	if !isFlagPassed("email") || !isFlagPassed("token") {
		flag.Usage()
		return
	}
	rotateproxy.StartRunCrawler(token, email, rule, pageCount)
	rotateproxy.StartCheckProxyAlive(timeout)
	c := rotateproxy.NewRedirectClient(rotateproxy.WithConfig(&baseCfg))
	c.Serve()
	select {}
}
