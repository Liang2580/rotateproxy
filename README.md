# rotateproxy

利用fofa搜索socks5开放代理进行代理池轮切的工具

## 帮助

```shell
└──╼ #./cmd/rotateproxy/rotateproxy -h
Usage of ./cmd/rotateproxy/rotateproxy:
  -email string
    	email address
  -l string
    	listen address (default ":9999")
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


使用方法

./cmd/rotateproxy/rotateproxy  -email=qqq.com  -token=xxx 

time 代表。毫秒值。默认为超过8000毫秒则不使用当前代理

![proxy](https://user-images.githubusercontent.com/27684409/142030959-588afe68-3a1a-4734-86e4-7ec347e18e21.png)




##此项目还欠缺一部分功能


1.当前代理中断后的自动连接耗时最短的代理  [2.0已完成]

2.调度的权重

3.暂只编译了Linux amd64 
