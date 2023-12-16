package main

import (
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/DalphanDev/Turbo/http"
	"github.com/DalphanDev/Turbo/http/cookiejar"
	"github.com/andybalholm/brotli"
)

type Field struct {
	Name   string
	Value  string
	Inline bool
}

type Footer struct {
	Text    string
	IconURL string
}

type Embed struct {
	Title  string
	Color  int
	Fields []Field
	Footer Footer
}

type WebhookPayload struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Embeds    []Embed
}

type ClientResponse struct {
	Command  string `json:"command"`
	ClientID string `json:"clientID"`
}

type DoResponse struct {
	Command    string      `json:"command"`
	ClientID   string      `json:"clientID"`
	StatusCode int         `json:"statusCode"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body"`
}

type TurboClient struct {
	Client *http.Client // Not exported because lowercase letters, which is fine, same for proxy
	Proxy  *url.URL
}

type RequestOptions struct {
	URL     string
	Headers map[string]string
	Body    io.Reader
}

type TurboResponse struct {
	StatusCode int
	Headers    http.Header
	Body       string
}

func NewTurboClient(proxy string, mimicSetting string) (*TurboClient, error) {
	transport := &http.Transport{
		IdleConnTimeout: 10000 * time.Millisecond,
		MimicSetting:    mimicSetting,
	}

	// Setup the http2 transport

	// // Configure transport for HTTP/2
	// http2Transport := &http2.Transport{
	// 	AllowHTTP:                 true, // if you want to allow non-TLS requests
	// 	MaxEncoderHeaderTableSize: 65536,
	// 	MaxDecoderHeaderTableSize: 65536,
	// }

	// // Configure HTTP/2 transport
	// http2Transport := &http2.Transport{
	// 	AllowHTTP:         true, // if you want to allow non-TLS requests
	// 	MaxHeaderListSize: 65536,
	// }

	// // This is crucial: wrap the standard transport with HTTP/2-specific settings
	// transport.RegisterProtocol("h2", http2Transport)

	// // Apply HTTP/2 settings to the custom transport
	// if err := http2.ConfigureTransport(transport); err != nil {
	// 	return nil, err
	// }

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: transport,
		Jar:       jar,
	}

	if proxy != "" {
		proxyURL, err := parseProxy(proxy)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxyURL)
		return &TurboClient{
			Client: client,
			Proxy:  proxyURL,
		}, nil
	}

	return &TurboClient{
		Client: client,
		Proxy:  nil,
	}, nil
}

// Do sends an HTTP request and returns an HTTP response.
func (tc *TurboClient) Do(method string, options RequestOptions) (*TurboResponse, error) {
	req, err := http.NewRequest(method, options.URL, options.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Request successfully created.")

	var defaultHeaders = map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	}

	// fmt.Println("HTTP Request Details:")
	// fmt.Println("URL:", options.URL)
	// fmt.Println("Method:", method)
	// fmt.Println("Headers:")

	mergedHeaders := mergeHeaders(defaultHeaders, options.Headers)

	// fmt.Println(mergedHeaders)

	for key, value := range mergedHeaders {
		// fmt.Println(key+":", value)
		req.Header.Set(key, value)
	}

	// Check if the proxy requires authentication
	if tc.Proxy != nil && tc.Proxy.User != nil {
		username := tc.Proxy.User.Username()
		password, _ := tc.Proxy.User.Password()
		if username != "" && password != "" {
			auth := username + ":" + password
			encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
			req.Header.Set("Proxy-Authorization", "Basic "+encodedAuth)
		}
	}

	fmt.Println("Sending request...")
	fmt.Println(req)

	resp, err := tc.Client.Do(req)
	if err != nil {
		fmt.Println("error in client.Do")
		fmt.Println(err)
		return nil, err
	}

	myResponse := &TurboResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gz.Close()
		body, err := ioutil.ReadAll(gz)
		if err != nil {
			return nil, err
		}
		// Use body for the decompressed response
		myResponse.Body = string(body)
	} else if resp.Header.Get("Content-Encoding") == "br" {
		br := brotli.NewReader(resp.Body)
		body, err := ioutil.ReadAll(br)
		if err != nil {
			return nil, err
		}
		myResponse.Body = string(body)

	} else if resp.Header.Get("Content-Encoding") == "deflate" {
		fl := flate.NewReader(resp.Body)
		defer fl.Close()
		body, err := ioutil.ReadAll(fl)
		if err != nil {
			return nil, err
		}
		myResponse.Body = string(body)

	} else {
		fmt.Println(resp.Header.Get("Content-Encoding"))
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		myResponse.Body = string(body)
	}

	// Now format the response into an object that we can use easily.
	return myResponse, nil
}

func mergeHeaders(defaultHeaders, customHeaders map[string]string) map[string]string {
	mergedHeaders := make(map[string]string)
	for key, value := range defaultHeaders {
		lowerKey := strings.ToLower(key)
		mergedHeaders[lowerKey] = value
	}
	for key, value := range customHeaders {
		lowerKey := strings.ToLower(key)
		mergedHeaders[lowerKey] = value
	}
	return mergedHeaders
}

func parseProxy(proxyString string) (*url.URL, error) {
	components := strings.Split(proxyString, ":")

	if len(components) != 4 {
		return nil, fmt.Errorf("invalid proxy string format")
	}

	ip := components[0]
	port := components[1]
	username := components[2]
	password := components[3]

	proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", username, password, ip, port)
	// proxyURL, err := url.Parse(proxyURL)
	parsedProxyURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}
	return parsedProxyURL, nil
}

func main() {
	// Example proxy: 207.90.213.151:15413:egvrca423:qhYCz8388o
	client, err := NewTurboClient("", "chrome")

	if err != nil {
		panic(err)
	}

	headers := map[string]string{
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36", // This will overwrite the default User-Agent header
		"accept":                    "application/json",
		"accept-language":           "en-US,en;q=0.9",
		"accept-encoding":           "gzip, deflate, br", // Test with no deflate header.
		"content-type":              "application/json",
		"dnt":                       "1",
		"sec-ch-ua":                 "\"Chromium\";v=\"116\", \"Not)A;Brand\";v=\"24\", \"Google Chrome\";v=\"116\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "Windows",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
	}

	// body := "your string data"

	options := RequestOptions{
		// URL:     "https://deposit.us.shopifycs.com/sessions",
		URL:     "https://cncpts.com/",
		Headers: headers,
		Body:    nil, // Can either use nil or a string reader.
	}

	resp, err := client.Do("GET", options)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
