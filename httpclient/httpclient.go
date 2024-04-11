package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MakeHTTPRequest хүсэлтийн тохиргоо
//
// Url: хүсэлтийн url заавал байх ёстой
//
// Method: хүсэлтийн төрөл явуулаагүй үед GET байна
//
// Headers: default-р Content-Type нь application/json байна
//
// Parameters: query parameters
//
// Body: хүсэлтийн бие
//
// Timeout: хүсэлтийн timeout явуулаагүй үед 10сек байна
type Config struct {
	Url        string
	Method     string
	Headers    *map[string]string
	Parameters *url.Values
	Body       []byte
	Timeout    uint
}

// Response
//
// IsSuccess: Амжилттай бол true байна.
//
// Code: 0 бол client кодын алдаа бусад тохиолдолд httpStatusCode байна.
//
// Message:
type Response struct {
	IsSuccess bool
	Code      int
	Message   string
	Body      []byte
}

// Send: send http request
func Send(config *Config) *Response {
	if config.Method == "" {
		config.Method = "GET"
	}

	if config.Timeout == 0 {
		config.Timeout = 10
	}

	client := http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	u, err := url.Parse(config.Url)
	if err != nil {
		return &Response{IsSuccess: false, Code: 0, Message: err.Error()}
	}

	config.Method = strings.ToUpper(config.Method)

	q := u.Query()
	if config.Parameters != nil {
		for k, v := range *config.Parameters {
			q.Set(k, strings.Join(v, ","))
		}
	}
	u.RawQuery = q.Encode()

	req := &http.Request{}
	if config.Body != nil {
		req, err = http.NewRequest(config.Method, u.String(), bytes.NewBuffer(config.Body))
		if err != nil {
			return &Response{IsSuccess: false, Code: 0, Message: err.Error()}
		}
	} else {
		req, err = http.NewRequest(config.Method, u.String(), nil)
		if err != nil {
			return &Response{IsSuccess: false, Code: 0, Message: err.Error()}
		}
	}

	if config.Headers != nil {
		for k, v := range *config.Headers {
			req.Header.Set(k, v)
		}
	}

	httpResponse, httpError := client.Do(req)
	if httpError != nil {
		return &Response{IsSuccess: false, Code: 0, Message: httpError.Error()}
	}

	if httpResponse == nil {
		return &Response{IsSuccess: false, Code: 0, Message: "httpResponse empty"}
	}

	log.Printf("%d %s %s\n", httpResponse.StatusCode, config.Method, req.URL.String())

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return &Response{IsSuccess: false, Code: 0, Message: err.Error()}
	}
	defer httpResponse.Body.Close()

	return &Response{IsSuccess: httpResponse.StatusCode == http.StatusOK,
		Code:    httpResponse.StatusCode,
		Message: httpResponse.Status,
		Body:    body}
}
