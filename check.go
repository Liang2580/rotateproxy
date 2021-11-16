package rotateproxy

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type IPResponse struct {
	ORIGIN    string     `json:"origin"`
}


type IPInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func CheckProxyAlive(proxyURL,ip string,timeout int) (respBody string, avail bool,time_config int) {
	startTime := time.Now().UnixNano()
	proxy, _ := url.Parse(proxyURL)
	httpclient := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 20 * time.Second,
	}
	resp, err := httpclient.Get("http://httpbin.org/ip")
	if err != nil {
		return "", false,10000
	}
	defer resp.Body.Close()
	var res IPResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if len(res.ORIGIN)<2{
		return "", false,10000
	}
	if err != nil {
		return "", false,10000
	}
	endTime := time.Now().UnixNano()
	Milliseconds:= int((endTime - startTime) / 1e6)// 毫秒
	fmt.Println("IP : ",ip,"  用时毫秒 ",Milliseconds,"   可用")
	if Milliseconds>timeout{
		return "", false,10000
	}
	return res.ORIGIN, true,Milliseconds
}

func StartCheckProxyAlive(timeout int) {
	go func() {
		ticker := time.NewTicker(120 * time.Second)
		for {
			select {
			case <-crawlDone:
				fmt.Println("Checking")
				checkAlive(timeout)
				fmt.Println("Check done")
			case <-ticker.C:
				checkAlive(timeout)
			}
		}
	}()
}

func checkAlive(timeout int) {
	proxies, err := QueryProxyURL()
	if err != nil {
		fmt.Printf("[!] query db error: %v\n", err)
	}
	for i := range proxies {
		proxy := proxies[i]
		if proxy.Available{
			continue
		}
		if proxy.Retry >5{
			continue
		}
		go func() {
			_, _,time_conig := CheckProxyAlive(proxy.URL,proxy.IP,timeout)
			if time_conig !=10000 {
				SetProxytime(proxy.URL, time_conig)
			}else{
				AddProxyURLRetry(proxy.URL)
			}
		}()
	}
}
