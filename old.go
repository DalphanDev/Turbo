package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	tls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

type turboOptions struct {
	Body      string
	UserAgent string
	Proxy     string
	Cookies   string
}

var dialTimeout = time.Duration(15) * time.Second
var requestHostname = "facebook.com" // speaks http2 and TLS 1.3
var requestAddr = "31.13.72.36:443"

func fart() {
	// var response *http.Response
	// var err error

	// config := tls.Config{ServerName: "google.com"}

	// tcpConn, err := tls.Dial("tcp", "142.250.217.164:443", &config)

	// if err != nil {
	// 	panic("failed to connect: " + err.Error())
	// }

	// // Create a tls client using config, spec.
	// uTlsConn := tls.UClient(tcpConn, &config, tls.HelloChrome_Auto)

	// defer uTlsConn.Close()

	// err = uTlsConn.Handshake()
	// if err != nil {
	// 	fmt.Printf("uTlsConn.Handshake() error: %+v", err)
	// }

	dialConn, err := net.Dial("tcp", "31.13.72.36:443")
	if err != nil {
		fmt.Printf("net.Dial() failed: %+v\n", err)
		return
	}

	config := tls.Config{ServerName: "facebook.com"}
	tlsConn := tls.UClient(dialConn, &config, tls.HelloChrome_62)
	// n, err := tlsConn.Write("Hello, World!")
	response, err := httpGetOverConn(tlsConn, tlsConn.HandshakeState.ServerHello.AlpnProtocol)
	if err != nil {
		fmt.Printf("#> HttpGetByHelloID(HelloChrome_62) failed: %+v\n", err)
	} else {
		fmt.Printf("#> HttpGetByHelloID(HelloChrome_62) response: %+s\n", dumpResponseNoBody(response))
	}

	// response, err := HttpGetByHelloID(requestHostname, requestAddr, tls.HelloChrome_Auto)
	// if err != nil {
	// 	fmt.Printf("#> HttpGetByHelloID(HelloChrome_62) failed: %+v\n", err)
	// } else {
	// 	fmt.Printf("#> HttpGetByHelloID(HelloChrome_62) response: %+s\n", dumpResponseNoBody(response))
	// }

}

func HttpGetCustom(targetURL string) (*http.Response, error) {

	// Parse the URL to get the hostname
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("url.Parse error: %+v", err)
	}

	addr, err := getIPv4Address(u.Hostname())
	if err != nil {
		return nil, err
	}

	fmt.Println("Address: " + addr)

	config := tls.Config{ServerName: u.Hostname()}

	dialConn, err := net.DialTimeout("tcp", addr, dialTimeout)
	if err != nil {
		return nil, fmt.Errorf("net.DialTimeout error: %+v", err)
	}
	uTlsConn := tls.UClient(dialConn, &config, tls.HelloCustom)
	defer uTlsConn.Close()

	// do not use this particular spec in production
	// make sure to generate a separate copy of ClientHelloSpec for every connection
	// spec := tls.ClientHelloSpec{
	// 	TLSVersMax: tls.VersionTLS13,
	// 	TLSVersMin: tls.VersionTLS10,
	// 	CipherSuites: []uint16{ // Set custom cipher suites to mimic google chrome
	// 		tls.GREASE_PLACEHOLDER,
	// 		tls.TLS_AES_128_GCM_SHA256,
	// 		tls.TLS_AES_256_GCM_SHA384,
	// 		tls.TLS_CHACHA20_POLY1305_SHA256,
	// 		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	// 		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
	// 		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	// 		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	// 		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	// 		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	// 	},
	// 	Extensions: []tls.TLSExtension{ // Set custom extensions to mimic google chrome
	// 		&tls.UtlsGREASEExtension{},
	// 		&tls.SessionTicketExtension{},
	// 		&tls.KeyShareExtension{[]tls.KeyShare{
	// 			{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
	// 			{Group: tls.X25519},
	// 		}},
	// 		&tls.StatusRequestExtension{},
	// 		&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
	// 		&tls.SCTExtension{},
	// 		&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
	// 		&tls.SupportedCurvesExtension{
	// 			Curves: []tls.CurveID{
	// 				tls.CurveID(tls.GREASE_PLACEHOLDER),
	// 				tls.X25519,
	// 				tls.CurveP256,
	// 				tls.CurveP384,
	// 			},
	// 		},
	// 		&tls.SupportedPointsExtension{SupportedPoints: []byte{0x00}},
	// 		&tls.PSKKeyExchangeModesExtension{
	// 			Modes: []uint8{
	// 				tls.PskModeDHE,
	// 			}},
	// 		&tls.SNIExtension{},
	// 		&tls.UtlsCompressCertExtension{
	// 			Algorithms: []tls.CertCompressionAlgo{
	// 				tls.CertCompressionBrotli,
	// 			}},
	// 		&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
	// 			tls.ECDSAWithP256AndSHA256,
	// 			tls.PSSWithSHA256,
	// 			tls.PKCS1WithSHA256,
	// 			tls.ECDSAWithP384AndSHA384,
	// 			tls.PSSWithSHA384,
	// 			tls.PKCS1WithSHA384,
	// 			tls.PSSWithSHA512,
	// 			tls.PKCS1WithSHA512,
	// 		}},
	// 		&tls.UtlsExtendedMasterSecretExtension{},
	// 		// Supported points might go here instead
	// 		&tls.SupportedVersionsExtension{Versions: []uint16{
	// 			tls.GREASE_PLACEHOLDER,
	// 			tls.VersionTLS13,
	// 			tls.VersionTLS12,
	// 		}},
	// 		&tls.ApplicationSettingsExtension{
	// 			SupportedProtocols: []string{
	// 				"h2",
	// 			},
	// 		},
	// 		&tls.UtlsGREASEExtension{},
	// 	},
	// }

	spec := tls.ClientHelloSpec{
		TLSVersMax: tls.VersionTLS13,
		TLSVersMin: tls.VersionTLS10,
		CipherSuites: []uint16{
			tls.GREASE_PLACEHOLDER,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Extensions: []tls.TLSExtension{
			&tls.UtlsGREASEExtension{},
			&tls.SessionTicketExtension{},
			&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
				{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
				{Group: tls.X25519},
			}},
			&tls.StatusRequestExtension{},
			&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
			&tls.SCTExtension{},
			&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
			&tls.SupportedCurvesExtension{
				Curves: []tls.CurveID{
					tls.CurveID(tls.GREASE_PLACEHOLDER),
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				},
			},
			&tls.SupportedPointsExtension{SupportedPoints: []byte{0x00}},
			&tls.PSKKeyExchangeModesExtension{
				Modes: []uint8{
					tls.PskModeDHE,
				}},
			&tls.SNIExtension{},
			&tls.UtlsCompressCertExtension{
				Algorithms: []tls.CertCompressionAlgo{
					tls.CertCompressionBrotli,
				}},
			&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
				tls.ECDSAWithP256AndSHA256,
				tls.PSSWithSHA256,
				tls.PKCS1WithSHA256,
				tls.ECDSAWithP384AndSHA384,
				tls.PSSWithSHA384,
				tls.PKCS1WithSHA384,
				tls.PSSWithSHA512,
				tls.PKCS1WithSHA512,
			}},
			&tls.UtlsExtendedMasterSecretExtension{},
			&tls.SupportedVersionsExtension{Versions: []uint16{
				tls.GREASE_PLACEHOLDER,
				tls.VersionTLS13,
				tls.VersionTLS12,
			}},
			&tls.ApplicationSettingsExtension{
				SupportedProtocols: []string{
					"h2",
				},
			},
		},
	}

	err = uTlsConn.ApplyPreset(&spec)

	if err != nil {
		return nil, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}

	err = uTlsConn.Handshake()
	if err != nil {
		return nil, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}

	return httpGetOverConn(uTlsConn, uTlsConn.HandshakeState.ServerHello.AlpnProtocol)
}

func httpGetOverConn(conn net.Conn, alpn string) (*http.Response, error) {
	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Host: "www." + requestHostname + "/"},
		Header: make(http.Header),
		Host:   "www." + requestHostname,
	}

	switch alpn {
	case "h2":
		req.Proto = "HTTP/2.0"
		req.ProtoMajor = 2
		req.ProtoMinor = 0

		tr := http2.Transport{}
		cConn, err := tr.NewClientConn(conn)
		if err != nil {
			return nil, err
		}
		return cConn.RoundTrip(req)
	case "http/1.1", "":
		req.Proto = "HTTP/1.1"
		req.ProtoMajor = 1
		req.ProtoMinor = 1

		err := req.Write(conn)
		if err != nil {
			return nil, err
		}
		return http.ReadResponse(bufio.NewReader(conn), req)
	default:
		return nil, fmt.Errorf("unsupported ALPN: %v", alpn)
	}
}

func dumpResponseNoBody(response *http.Response) string {
	resp, err := httputil.DumpResponse(response, false)
	if err != nil {
		return fmt.Sprintf("failed to dump response: %v", err)
	}
	return string(resp)
}

func getIPv4Address(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return "", fmt.Errorf("failed to lookup IP address for hostname %s: %v", hostname, err)
	}

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return net.JoinHostPort(ipv4.String(), "443"), nil
		}
	}

	return "", fmt.Errorf("no IPv4 address found for hostname %s", hostname)
}

func HttpGetByHelloID(hostname string, addr string, helloID tls.ClientHelloID) (*http.Response, error) {
	config := tls.Config{ServerName: hostname}
	dialConn, err := net.DialTimeout("tcp", addr, dialTimeout)
	if err != nil {
		return nil, fmt.Errorf("net.DialTimeout error: %+v", err)
	}
	uTlsConn := tls.UClient(dialConn, &config, helloID)
	defer uTlsConn.Close()

	err = uTlsConn.Handshake()
	if err != nil {
		return nil, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}

	return httpGetOverConn(uTlsConn, uTlsConn.HandshakeState.ServerHello.AlpnProtocol)
}
