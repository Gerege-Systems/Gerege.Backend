package httpclient

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
// Timeout: хүсэлтийн timeout явуулаагүй үед 45 сек байна
type HttpConfig struct {
	Url        string
	Method     string
	Headers    *map[string]string
	Parameters *url.Values
	Body       []byte
	Timeout    uint
}

// HttpResult
//
// IsSuccess: Амжилттай бол true байна.
//
// Code: 0 бол client кодын алдаа бусад тохиолдолд httpStatusCode байна.
//
// Message:
type HttpResult struct {
	IsSuccess bool
	Code      int
	Message   string
	Body      []byte
}

// Send: send http request
func Send(httpConfig *HttpConfig) *HttpResult {
	if httpConfig.Method == "" {
		httpConfig.Method = "GET"
	}
	httpConfig.Method = strings.ToUpper(httpConfig.Method)

	if httpConfig.Timeout == 0 {
		httpConfig.Timeout = 45
	}

	requestUrl, err := url.Parse(httpConfig.Url)
	if err != nil {
		return &HttpResult{IsSuccess: false, Code: 0, Message: err.Error()}
	}

	urlValues := requestUrl.Query()
	if httpConfig.Parameters != nil {
		for k, v := range *httpConfig.Parameters {
			urlValues.Set(k, strings.Join(v, ","))
		}
	}
	requestUrl.RawQuery = urlValues.Encode()

	httpRequest, err := http.NewRequest(httpConfig.Method, requestUrl.String(), bytes.NewBuffer(httpConfig.Body))
	if err != nil {
		return &HttpResult{IsSuccess: false, Code: 0, Message: err.Error()}
	}

	if httpConfig.Headers != nil {
		for k, v := range *httpConfig.Headers {
			httpRequest.Header.Set(k, v)
		}
	}

	httpClient := http.Client{
		Timeout: time.Duration(httpConfig.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	httpResponse, httpError := httpClient.Do(httpRequest)
	if httpError != nil {
		return &HttpResult{IsSuccess: false, Code: 0, Message: httpError.Error()}
	}

	if httpResponse == nil {
		return &HttpResult{IsSuccess: false, Code: 0, Message: "httpResponse empty"}
	}

	log.Printf("%d %s %s\n", httpResponse.StatusCode, httpConfig.Method, httpRequest.URL.String())

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return &HttpResult{IsSuccess: false, Code: 0, Message: err.Error()}
	}
	defer httpResponse.Body.Close()

	return &HttpResult{IsSuccess: httpResponse.StatusCode == http.StatusOK,
		Code:    httpResponse.StatusCode,
		Message: httpResponse.Status,
		Body:    body}
}
