// package main

// import (
// 	// "compress/gzip"
// 	"fmt"
// 	"io/ioutil"
// 	"net"
// 	"net/url"
// 	"time"

// 	"github.com/DalphanDev/Turbo/http"

// 	// "github.com/DalphanDev/Turbo/mimic"
// 	tls "github.com/refraction-networking/utls"
// )

// func main() {
// 	// Test a turbo request!
// 	//
// 	// What are the steps to making a request with uTLS?

// 	targetURL := "https://example.com/"

// 	parsedURL, err := url.Parse(targetURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	serverName := parsedURL.Host

// 	fmt.Println(serverName)

// 	// Make an http transport using our custom Dial TLS function.

// 	transport := &http.Transport{
// 		// DialTLS: DialWithUTLS, // Comment this out to test uTLS vs native TLS
// 	}

// 	client := &http.Client{
// 		Transport: transport,
// 	}

// 	resp, err := client.Get(targetURL)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(string(body))
// }

// func DialWithUTLS(network, addr string) (*tls.UConn, error) {
// 	// create a dialer object
// 	dialer := &net.Dialer{
// 		Timeout:   time.Second * 30,
// 		KeepAlive: time.Second * 30,
// 		DualStack: true,
// 	}

// 	// establish a TCP connection to the remote server
// 	conn, err := dialer.Dial(network, addr)
// 	if err != nil {
// 		fmt.Println("TCP Connection Failed!")
// 	}

// 	chromeAuto := tls.HelloChrome_58

// 	tlsConn := tls.UClient(conn, &tls.Config{
// 		ServerName:         addr,
// 		InsecureSkipVerify: true,
// 	}, chromeAuto)

// 	if err != nil {
// 		fmt.Printf("uTLSConn generation error: %+v", err)
// 	}

// 	// perform the uTLS handshake
// 	err = tlsConn.Handshake()
// 	if err != nil {
// 		conn.Close()
// 		fmt.Println("TLS Handshake Failed!")
// 	}

// 	return tlsConn, nil
// }
