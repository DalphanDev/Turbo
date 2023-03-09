package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"time"

	"github.com/DalphanDev/Turbo/http"
	"github.com/DalphanDev/Turbo/http/httptrace"
	"github.com/DalphanDev/Turbo/mimic"
	tls "github.com/refraction-networking/utls"
)

// type turboOptions struct {
// 	Body      string
// 	UserAgent string
// 	Proxy     string
// 	Cookies   string
// }

func main() {
	// Test a turbo request!

	// What are the steps to making a request with uTLS?

	// Well first, let's try and fetch a preset uTLS fingerprint.
	modernChrome := mimic.NewChromeMimic.ClientHello()

	// targetAddress := "example.com:443"

	// targetURL := "https://eoobxe7m89qj9cl.m.pipedream.net"
	// targetURL := "https://purple.com/"
	// targetURL := "https://www.google.com/"
	// targetURL := "https://www.whatsmybrowser.org/"
	targetURL := "https://example.com/"

	// TEST API
	// targetAddress := "https://eoobxe7m89qj9cl.m.pipedream.net"

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	serverName := parsedURL.Host

	fmt.Println(serverName)

	targetAddress := net.JoinHostPort(serverName, "443")

	fmt.Println(targetAddress)

	// The first step in making a request to any server, is establishing a TCP connection.
	tcpConn, err := net.Dial("tcp", targetAddress)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Successfully established a TCP connection to %s\n", targetAddress)
	}

	// The next step in making a request to any server, is to send the Client Hello Message.

	// We need a couple things to do this. First we need a TCP connection, which we have already created above.

	// Next, we need a tls config which represents the TLS configuration used by a TLS client or server.
	config := tls.Config{
		ServerName: serverName,
	}

	// Lastly, we need a clientHelloID to pass into our UClient. Something like tls.HelloChrome_62
	// However, since we are using a custom client hello spec, we need to use tls.HelloCustom

	uTlsConn := tls.UClient(tcpConn, &config, tls.HelloCustom)

	defer uTlsConn.Close()
	err = uTlsConn.ApplyPreset(modernChrome)

	fmt.Println(uTlsConn)

	if err != nil {
		fmt.Printf("uTlsConn.Handshake() error: %+v", err)
		return
	}

	err = uTlsConn.Handshake()

	if err != nil {
		fmt.Printf("uTlsConn.Handshake() error: %+v", err)
		return
	}

	fmt.Printf("Successful handshake made to %s\n", targetAddress)
	fmt.Printf("ALPN PROTOCOL: %s\n", uTlsConn.HandshakeState.ServerHello.AlpnProtocol)
	alpn := uTlsConn.HandshakeState.ServerHello.AlpnProtocol
	// FYI, not all servers return an ALPN Protocol. If it is empty, default to HTTP 1.1

	// The next step is to send an HTTP request over the established TLS connection.

	// First, create the HTTP request
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Printf("http.NewRequest() error: %+v\n", err)
		return
	}

	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("accept-encoding", "gzip, deflate, br")
	req.Header.Add("dnt", "1")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	trace := &httptrace.ClientTrace{
        TLSHandshakeStart: func() {
            fmt.Println("TLS handshake started")
        },
        TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
            fmt.Println("TLS handshake done")
            if err != nil {
                fmt.Println("TLS handshake error:", err)
            } else {
                fmt.Println("TLS handshake cipher suite:", cs.CipherSuite)
                fmt.Println("TLS handshake version:", cs.Version)
                fmt.Println("TLS handshake verified chains:", len(cs.VerifiedChains))
            }
        },
    }

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	fmt.Println(req)

	fmt.Println("ALPN PROTOCOL: ", alpn)

	// var resp *http.Response

	switch alpn {
	case "h2":
		req.Proto = "HTTP/2.0"
		req.ProtoMajor = 2
		req.ProtoMinor = 0

		fmt.Println("HTTP/2.0")

		// tr := http2.Transport{}
		// cConn, err := tr.NewClientConn(uTlsConn)
		// if err != nil {
		// 	fmt.Println("Error writing HTTP 2")
		// 	return
		// }
		// resp, err = cConn.RoundTrip(req)
		// if err != nil {
		// 	fmt.Println("Error reading HTTP 2")
		// 	return
		// }
	case "http/1.1", "":
		req.Proto = "HTTP/1.1"
		req.ProtoMajor = 1
		req.ProtoMinor = 1

		fmt.Println("HTTP/1.1")

		// create the transport
		transport := &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			DialTLS: DialWithUTLS, // Comment this out to test uTLS vs native TLS
		}

		fmt.Println("uTLS transport created!")

		client := &http.Client{
			Transport: transport,
		}

		fmt.Println("uTLS client created!")
		
		resp, err := client.Do(req)
		if err != nil {
			// handle error
			fmt.Println("Error sending HTTP 1.1 request...")
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		fmt.Println(resp.Header.Get("Content-Encoding"))

		if resp.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(resp.Body)
			if err != nil {
				// handle error
			}
			defer gz.Close()
			body, err := ioutil.ReadAll(gz)
			if err != nil {
				// handle error
			}
			// Use body for the decompressed response

			fmt.Println(string(body))
		}

	default:
		fmt.Errorf("unsupported ALPN: %v", alpn)
		return
	}
}

func DialWithUTLS(network, addr string) (net.Conn, error) {

	// fmt.Println("DialWithUTLS Called!")

    // create a dialer object
    dialer := &net.Dialer{
        Timeout:   time.Second * 30,
        KeepAlive: time.Second * 30,
        DualStack: true,
    }

	// fmt.Println("Dialer Created!")

    // establish a TCP connection to the remote server
    conn, err := dialer.Dial(network, addr)
    if err != nil {
		fmt.Println("TCP Connection Failed!")
		return nil, err
    }

	// fmt.Println("TCP Connection Established!")

	modernChrome := mimic.NewChromeMimic.ClientHello()

	// fmt.Println("modernChrome fingerprint fetched!")

	tlsConn := tls.UClient(conn, &tls.Config{
		ServerName:         addr,
		InsecureSkipVerify: true,
	}, tls.HelloCustom)

	// fmt.Println("tlsConn created!")

	// defer tlsConn.Close() We use this to close our connection after our request is complete.
	err = tlsConn.ApplyPreset(modernChrome)

	if err != nil {
		fmt.Printf("uTLSConn generation error: %+v", err)
	}

    // perform the uTLS handshake
    err = tlsConn.Handshake()
    if err != nil {
        conn.Close()
        fmt.Println("TLS Handshake Failed!")
    }

	// fmt.Println("TLS Handshake Completed!")

	// fmt.Println("Returning TLS Connection!")

    return tlsConn.Conn, nil
}
