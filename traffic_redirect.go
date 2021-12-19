package rotateproxy

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

var (
	largeBufferSize = 32 * 1024 // 32KB large buffer
)

type BaseConfig struct {
	ListenAddr   string
	IPRegionFlag int // 0: all 1: cannot bypass gfw 2: bypass gfw
}

type RedirectClient struct {
	config *BaseConfig
}

type RedirectClientOption func(*RedirectClient)

func WithConfig(config *BaseConfig) RedirectClientOption {
	return func(c *RedirectClient) {
		c.config = config
	}
}

func NewRedirectClient(opts ...RedirectClientOption) *RedirectClient {
	c := &RedirectClient{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *RedirectClient) Serve() error {
	l, err := net.Listen("tcp", c.config.ListenAddr)
	if err != nil {
		return err
	}
	//for IsProxyURLBlank() {
	//	fmt.Println("[*] waiting for crawl proxy...")
	//	time.Sleep(3 * time.Second)
	//}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("[!] accept error: %v\n", err)
			continue
		}
		go c.HandleConn(conn)
	}
}


func (c *RedirectClient) HandleConn(conn net.Conn) {
	startTime := time.Now().UnixNano()
	key, err := RandomProxyURL(c.config.IPRegionFlag)
	if err != nil {
		errConn := closeConn(conn)
		if errConn != nil {
			fmt.Printf("[!] close connect error : %v\n", errConn)
		}
		return
	}
	key2 := strings.TrimPrefix(key, "socks5://")
	fmt.Println(key)
	cc, err := net.DialTimeout("tcp", key2, 500*time.Millisecond)
	flag :=false
	if err != nil {
		//如果超时。则踢出这个节点
		fmt.Printf("连接失败\n")
		flag=true
		go StopProxy(key)
		go closeConn(conn)
	}
	if flag{
		key2 = strings.TrimPrefix(key, "socks5://")
		fmt.Printf("[!] 连接代理失败了。重新选择节点中 %v\n", key2)
		cc, err = net.DialTimeout("tcp", key2, 500*time.Millisecond)
		if err != nil {
			StopProxy(key)
			closeConn(conn)
			return
		}
	}

	endTime := time.Now().UnixNano()
	Milliseconds:= int((endTime - startTime) / 1e6)// 毫秒
	if Milliseconds>30000{
		fmt.Printf("[!] cannot connect to error2 %v\n", key2)
		StopProxy(key)
		closeConn(conn)
		return
	}
	go func() {
		err = transport(conn, cc)
		if err != nil {
			fmt.Printf("[!] connect error: %v\n", err)
			errConn := closeConn(conn)
			if errConn != nil {
				fmt.Printf("[!] close connect error: %v\n", errConn)
			}
			errConn = closeConn(cc)
			if errConn != nil {
				fmt.Printf("[!] close upstream connect error: %v\n", errConn)
			}
		}
	}()

}

func (c *RedirectClient) HandleConn2(conn net.Conn) {
	key, err := RandomProxyURL(c.config.IPRegionFlag)
	if err != nil {
		errConn := closeConn(conn)
		if errConn != nil {
			fmt.Printf("[!] close connect error: %v\n", errConn)
		}
		return
	}
	key2 := strings.TrimPrefix(key, "socks5://")
	fmt.Println(key)
	cc, err := net.DialTimeout("tcp", key2, 10*time.Second)
	flag:=true
	if err != nil {
		//如果超时。则踢出这个节点
		StopProxy(key)
		flag=flag
		fmt.Printf("[!] cannot connect to error %v\n", key2)
	}
	if !flag{
		url2, err := RandomProxyURL2()
		if err != nil {
			errConn := closeConn(conn)
			if errConn != nil {
				fmt.Printf("[!] close connect error2: %v\n", errConn)
			}
			return
		}
		key2 := strings.TrimPrefix(url2, "socks5://")
		fmt.Println(key)
		cc, err := net.DialTimeout("tcp", key2, 5*time.Second)
		if err != nil {
			//如果超时。则踢出这个节点
			StopProxy(key)
			flag=flag
			fmt.Printf("[!] cannot connect111 to %v\n", key2)
		}
		go func() {
			err = transport(conn, cc)
			if err != nil {
				fmt.Printf("[!] connect error222: %v\n", err)
				errConn := closeConn(conn)
				if errConn != nil {
					fmt.Printf("[!] close connect error333: %v\n", errConn)
				}
				errConn = closeConn(cc)
				if errConn != nil {
					fmt.Printf("[!] close upstream connect error44: %v\n", errConn)
				}
			}
		}()

	}else{
		go func() {
			err = transport(conn, cc)
			if err != nil {
				fmt.Printf("[!] connect error: %v\n", err)
				errConn := closeConn(conn)
				if errConn != nil {
					fmt.Printf("[!] close connect error: %v\n", errConn)
				}
				errConn = closeConn(cc)
				if errConn != nil {
					fmt.Printf("[!] close upstream connect error: %v\n", errConn)
				}
			}
		}()
	}
}

func closeConn(conn net.Conn) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	err = conn.Close()
	return err
}

func transport(rw1, rw2 io.ReadWriter) error {
	errc := make(chan error, 1)
	go func() {
		errc <- copyBuffer(rw1, rw2)
	}()

	go func() {
		errc <- copyBuffer(rw2, rw1)
	}()

	err := <-errc
	if err != nil && err == io.EOF {
		err = nil
	}
	return err
}

func copyBuffer(dst io.Writer, src io.Reader) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	buf := make([]byte, largeBufferSize)

	_, err = io.CopyBuffer(dst, src, buf)
	return err
}
