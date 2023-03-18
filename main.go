package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/DalphanDev/Turbo/http"
	"github.com/andybalholm/brotli"
)

type TurboClient struct {
	client *http.Client
	proxy  *url.URL
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

func NewTurboClient(proxy string) *TurboClient {
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
	}

	if proxy != "" {
		proxyURL, err := parseProxy(proxy)
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
		return &TurboClient{
			client: client,
			proxy:  proxyURL,
		}
	}

	return &TurboClient{
		client: client,
		proxy:  nil,
	}
}

// Do sends an HTTP request and returns an HTTP response.
func (tc *TurboClient) Do(method string, options RequestOptions) (*TurboResponse, error) {
	req, err := http.NewRequest(method, options.URL, options.Body)
	if err != nil {
		return nil, err
	}

	var defaultHeaders = map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language":           "en-US,en;q=0.9",
		"Accept-Encoding":           "gzip, deflate, br",
		"dnt":                       "1",
		"sec-ch-ua":                 "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"Windows\"",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	}

	mergedHeaders := mergeHeaders(defaultHeaders, options.Headers)
	for key, value := range mergedHeaders {
		req.Header.Set(key, value)
	}

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	// if tc.proxyURL != nil && tc.proxyUsername != "" && tc.proxyPassword != "" {
	// 	auth := tc.proxyUsername + ":" + tc.proxyPassword
	// 	fmt.Println(auth)
	// 	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	// 	fmt.Println(encodedAuth)
	// 	req.Header.Set("Proxy-Authenticate", "Basic "+encodedAuth)
	// }

	resp, err := tc.client.Do(req)
	if err != nil {
		panic(err)
	}

	myResponse := &TurboResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			// handle error
			fmt.Println("Error reading gzip: ", err)
		}
		defer gz.Close()
		body, err := ioutil.ReadAll(gz)
		if err != nil {
			// handle error
			fmt.Println("Error2 reading gzip: ", err)
		}
		// Use body for the decompressed response
		myResponse.Body = string(body)
	} else if resp.Header.Get("Content-Encoding") == "br" {
		br := brotli.NewReader(resp.Body)
		body, err := ioutil.ReadAll(br)
		if err != nil {
			panic(err)
		}
		myResponse.Body = string(body)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		myResponse.Body = string(body)
	}

	// Now format the response into an object that we can use easily.
	return myResponse, nil
}

func mergeHeaders(defaultHeaders, customHeaders map[string]string) map[string]string {
	mergedHeaders := make(map[string]string)
	for key, value := range defaultHeaders {
		mergedHeaders[key] = value
	}
	for key, value := range customHeaders {
		mergedHeaders[key] = value
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
		panic(err)
	}
	return parsedProxyURL, nil
}

func main() {
	// Example proxy: 207.90.213.151:15413:egvrca423:qhYCz8388o
	client := NewTurboClient("207.90.213.151:15413:egvrca423:qhYCz8388o")

	headers := map[string]string{
		"User-Agent": "Custom User Agent", // This will overwrite the default User-Agent header
	}

	body := "your string data"

	options := RequestOptions{
		URL:     "https://eoobxe7m89qj9cl.m.pipedream.net",
		Headers: headers,
		Body:    strings.NewReader(body), // Can either use nil or a string reader.
	}

	resp, err := client.Do("POST", options)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
