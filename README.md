# rotateproxy

利用fofa搜索socks5开放代理进行代理池轮切的工具

## 帮助

```shell
└──╼ #./cmd/rotateproxy/rotateproxy  -h
Usage of ./cmd/rotateproxy/rotateproxy:
  -email string
    	email address
  -l string
    	listen address (default ":9")
  -page int
    	the page count you want to crawl (default 5)
  -region int
    	0: all 1: cannot bypass gfw 2: bypass gfw
  -rule string
    	search rule (default "protocol==\"socks5\" && \"Version:5 Method:No Authentication(0x00)\" && after=\"2021-11-01\"")
  -time int
    	Not used for more than milliseconds (default 8000)
  -token string
    	token

```





