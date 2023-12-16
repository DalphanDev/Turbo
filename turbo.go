package turbo

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

	var defaultHeaders = map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		// "dnt":                       "1",
		// "sec-ch-ua":                 "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"",
		// "sec-ch-ua-mobile":          "?0",
		// "sec-ch-ua-platform":        "\"Windows\"",
		// "sec-fetch-dest":            "document",
		// "sec-fetch-mode":            "navigate",
		// "sec-fetch-site":            "none",
		// "sec-fetch-user":            "?1",
		// "upgrade-insecure-requests": "1",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
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

	// for key, value := range options.Headers {
	// 	req.Header.Set(key, value)
	// }

	// if method == "POST" || method == "PUT" {
	// 	bodyBytes, err := ioutil.ReadAll(options.Body)
	// 	if err != nil {
	// 		fmt.Println("Error reading request body:", err)
	// 		return nil, err
	// 	}
	// 	// Restore the body so it can be read again when the request is sent
	// 	options.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// 	fmt.Println("Request Body:", string(bodyBytes))
	// }

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
