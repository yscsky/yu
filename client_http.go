package yu

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HttpClient 嵌入http.Client，添加一些方法
type HttpClient struct {
	*http.Client
}

// NewHttpClient 创建HttpClient
func NewHttpClient() *HttpClient {
	return &HttpClient{
		Client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   15 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				ResponseHeaderTimeout: 30 * time.Second,
				ExpectContinueTimeout: 3 * time.Second,
				MaxIdleConns:          50,
				IdleConnTimeout:       60 * time.Second,
			},
		},
	}
}

// GetJSON GET请求，解析JSON
func (c *HttpClient) GetJSON(url string, data interface{}) (err error) {
	resp, err := c.Get(url)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return
}

// GetJSONAuth GET请求带BasicAuth认证，解析JSON
func (c *HttpClient) GetJSONAuth(url, un, pa string, data interface{}) (err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(un, pa)
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return
}

// PostJSON POST JSON请求，解析JSON
func (c *HttpClient) PostJSON(url string, src, data interface{}) (err error) {
	buffer := bytes.NewBuffer([]byte{})
	if err = json.NewEncoder(buffer).Encode(src); err != nil {
		return
	}
	resp, err := c.Post(url, "application/json;charset=utf-8", buffer)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return
}

// PostJSONAuth POST JSON请求带BasicAuth认证，解析JSON
func (c *HttpClient) PostJSONAuth(url, un, pa string, src, data interface{}) (err error) {
	buffer := bytes.NewBuffer([]byte{})
	if err = json.NewEncoder(buffer).Encode(src); err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, url, buffer)
	if err != nil {
		return
	}
	req.SetBasicAuth(un, pa)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return
}

// PostFormJSON POST请求，参数Form表单，解析JSON
func (c *HttpClient) PostFormJSON(url string, vals url.Values, data interface{}) (err error) {
	resp, err := c.PostForm(url, vals)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return
}

// PostFormJSONAuth POST请求带BasicAuth认证，参数Form表单，解析JSON
func (c *HttpClient) PostFormJSONAuth(url, un, pa string, vals url.Values, data interface{}) (err error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(vals.Encode()))
	if err != nil {
		return
	}
	req.SetBasicAuth(un, pa)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return
}
