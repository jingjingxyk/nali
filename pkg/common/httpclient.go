package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"

type HttpClient struct {
	*http.Client
}

var httpClient *HttpClient

const (
	EnableHttpProxy = true
	HttpProxyUrl    = "http://127.0.0.1:8015"
)

func init() {
	ProxyURL, _ := url.Parse(HttpProxyUrl)
	fmt.Println(http.ProxyURL(ProxyURL))
	fmt.Println(http.ProxyFromEnvironment)
	httpClient = &HttpClient{http.DefaultClient}
	httpClient.Timeout = time.Second * 60
	httpClient.Transport = &http.Transport{
		TLSHandshakeTimeout:   time.Second * 5,
		IdleConnTimeout:       time.Second * 10,
		ResponseHeaderTimeout: time.Second * 10,
		ExpectContinueTimeout: time.Second * 20,
		//Proxy:                 http.ProxyURL(ProxyURL),
		Proxy: http.ProxyFromEnvironment,
	}
}

func GetHttpClient() *HttpClient {
	c := *httpClient
	return &c
}

func (c *HttpClient) Get(urls ...string) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response

	for _, url := range urls {
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Println(err)
			continue
		}
		req.Header.Set("User-Agent", UserAgent)
		resp, err = c.Do(req)

		if err == nil && resp != nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			return
		}
	}

	return nil, err
}
