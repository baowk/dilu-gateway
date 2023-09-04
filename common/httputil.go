package common

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// HTTPClient 封装了HTTP客户端的功能
type HTTPClient struct {
	BaseURL string
	Headers map[string]string
}

func (c *HTTPClient) AddHeader(k, v string) {
	if c.Headers == nil {
		c.Headers = map[string]string{}
	}
	c.Headers[k] = v
}

// Get 发送GET请求
func (c *HTTPClient) Get(endpoint string) ([]byte, error) {
	url := c.BaseURL + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	return do(req)
}

// Post 发送POST请求
func (c *HTTPClient) Post(endpoint string, data []byte) ([]byte, error) {
	url := c.BaseURL + endpoint

	fmt.Println(url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	return do(req)
}

func do(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	fmt.Printf("%s : %s ", req.Method, req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
