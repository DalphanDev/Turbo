// package main

// import (
// 	"compress/gzip"
// 	"fmt"
// 	"io/ioutil"
// 	"net"
// 	"net/url"

// 	"github.com/DalphanDev/Turbo/http"
// 	"github.com/andybalholm/brotli"
// )

// func main() {
// 	// Test a turbo request!

// 	// What are the steps to making a request with uTLS?

// 	// targetURL := "https://example.com/"  ✅
// 	// targetURL := "https://www.whatsmybrowser.org/"  ✅
// 	// targetURL := "https://twitter.com/home" ✅
// 	// targetURL := "https://kith.com" ✅
// 	targetURL := "https://cncpts.com" // ✅ Access Granted. Had to edit the http2 transport to mimic google chrome!
// 	// targetURL := "https://tools.scrapfly.io/api/fp/ja3?extended=1" // Use this to see what your ja3 is!

// 	// Chrome fingerprint: [772,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,65281-0-16-27-11-18-13-5-23-43-10-45-17513-51-35,29-23-24,0]
// 	// Spoof fingerprint:  [772,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,45-51-17513-10-65281-5-43-23-0-16-11-35-13-27-18,29-23-24,0]
// 	parsedURL, err := url.Parse(targetURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	serverName := parsedURL.Host

// 	fmt.Println(serverName)

// 	targetAddress := net.JoinHostPort(serverName, "443")

// 	fmt.Println(targetAddress)

// 	transport := &http.Transport{}

// 	client := &http.Client{
// 		Transport: transport,
// 	}

// 	// First, create the HTTP request
// 	req, err := http.NewRequest("GET", targetURL, nil)
// 	if err != nil {
// 		fmt.Printf("http.NewRequest() error: %+v\n", err)
// 		return
// 	}

// 	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
// 	req.Header.Add("accept-language", "en-US,en;q=0.9")
// 	req.Header.Add("accept-encoding", "gzip, deflate, br")
// 	req.Header.Add("dnt", "1")
// 	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
// 	req.Header.Add("sec-ch-ua-mobile", "?0")
// 	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
// 	req.Header.Add("sec-fetch-dest", "document")
// 	req.Header.Add("sec-fetch-mode", "navigate")
// 	req.Header.Add("sec-fetch-site", "none")
// 	req.Header.Add("sec-fetch-user", "?1")
// 	req.Header.Add("upgrade-insecure-requests", "1")
// 	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

// 	// resp, err := client.Get(targetURL)
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	// body, err := ioutil.ReadAll(resp.Body)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	fmt.Println(resp.Header.Get("Content-Encoding")) // Encoding type

// 	if resp.Header.Get("Content-Encoding") == "gzip" {
// 		gz, err := gzip.NewReader(resp.Body)
// 		if err != nil {
// 			// handle error
// 			fmt.Println("Error reading gzip: ", err)
// 		}
// 		defer gz.Close()
// 		body, err := ioutil.ReadAll(gz)
// 		if err != nil {
// 			// handle error
// 			fmt.Println("Error2 reading gzip: ", err)
// 		}
// 		// Use body for the decompressed response

// 		fmt.Println(string(body))
// 	} else if resp.Header.Get("Content-Encoding") == "br" {
// 		br := brotli.NewReader(resp.Body)
// 		bodyBytes, err := ioutil.ReadAll(br)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(string(bodyBytes))
// 	} else {
// 		bytes, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(string(bytes))
// 	}

// 	fmt.Println("Request Status Code: ", resp.StatusCode)
// }
