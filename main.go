// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"net"
// 	"net/url"

// 	"github.com/DalphanDev/Turbo/http"
// 	"github.com/DalphanDev/Turbo/mimic"
// 	tls "github.com/refraction-networking/utls"
// 	"golang.org/x/net/http2"
// )

// // type turboOptions struct {
// // 	Body      string
// // 	UserAgent string
// // 	Proxy     string
// // 	Cookies   string
// // }

// func old_main() {
// 	// Test a turbo request!

// 	// What are the steps to making a request with uTLS?

// 	// Well first, let's try and fetch a preset uTLS fingerprint.
// 	// fingerprint := tls.HelloChrome_62
// 	modernChrome := mimic.NewChromeMimic.ClientHello()

// 	// targetAddress := "example.com:443"

// 	// targetURL := "https://eoobxe7m89qj9cl.m.pipedream.net"
// 	// targetURL := "https://purple.com/"
// 	// targetURL := "https://www.google.com/"
// 	// targetURL := "https://www.whatsmybrowser.org/"
// 	targetURL := "https://example.com/"

// 	// TEST API
// 	// targetAddress := "https://eoobxe7m89qj9cl.m.pipedream.net"

// 	parsedURL, err := url.Parse(targetURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	serverName := parsedURL.Host

// 	fmt.Println(serverName)

// 	targetAddress := net.JoinHostPort(serverName, "443")

// 	fmt.Println(targetAddress)

// 	// The first step in making a request to any server, is establishing a TCP connection.
// 	tcpConn, err := net.Dial("tcp", targetAddress)
// 	if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Printf("Successfully established a TCP connection to %s\n", targetAddress)
// 	}

// 	// The next step in making a request to any server, is to send the Client Hello Message.

// 	// We need a couple things to do this. First we need a TCP connection, which we have already created above.

// 	// Next, we need a tls config which represents the TLS configuration used by a TLS client or server.
// 	config := tls.Config{
// 		ServerName: serverName,
// 	}

// 	// Lastly, we need a clientHelloID to pass into our UClient. Something like tls.HelloChrome_62
// 	// However, since we are using a custom client hello spec, we need to use tls.HelloCustom

// 	uTlsConn := tls.UClient(tcpConn, &config, tls.HelloCustom)

// 	defer uTlsConn.Close()
// 	err = uTlsConn.ApplyPreset(modernChrome)

// 	if err != nil {
// 		fmt.Printf("uTlsConn.Handshake() error: %+v", err)
// 		return
// 	}

// 	err = uTlsConn.Handshake()

// 	if err != nil {
// 		fmt.Printf("uTlsConn.Handshake() error: %+v", err)
// 		return
// 	}

// 	fmt.Printf("Successful handshake made to %s\n", targetAddress)
// 	fmt.Printf("ALPN PROTOCOL: %s\n", uTlsConn.HandshakeState.ServerHello.AlpnProtocol) // For some reason doesn't print the ALPN protocol.
// 	alpn := uTlsConn.HandshakeState.ServerHello.AlpnProtocol
// 	// ^ Looking further, not all servers return an ALPN Protocol. If it is empty, default to HTTP 1.1

// 	// The next step is to send an HTTP request over the established TLS connection.

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

// 	fmt.Println(req)

// 	var resp *http.Response

// 	switch alpn {
// 	case "h2":
// 		req.Proto = "HTTP/2.0"
// 		req.ProtoMajor = 2
// 		req.ProtoMinor = 0

// 		tr := http2.Transport{}
// 		cConn, err := tr.NewClientConn(uTlsConn)
// 		if err != nil {
// 			fmt.Println("Error writing HTTP 2")
// 			return
// 		}
// 		resp, err = cConn.RoundTrip(req)
// 		if err != nil {
// 			fmt.Println("Error reading HTTP 2")
// 			return
// 		}
// 	case "http/1.1", "":
// 		req.Proto = "HTTP/1.1"
// 		req.ProtoMajor = 1
// 		req.ProtoMinor = 1

// 		// create the transport

// 	default:
// 		fmt.Errorf("unsupported ALPN: %v", alpn)
// 		return
// 	}

// 	fmt.Printf("Successful response: %d\n", resp.StatusCode)
// 	bodyBytes, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Printf("Error reading response body: %+v\n", err)
// 		return
// 	}
// 	fmt.Printf("Response Body: %s\n", string(bodyBytes))
// }

// func turbo() {

// }
