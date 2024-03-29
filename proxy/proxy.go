package proxy

import (
	"net"
	"github.com/gorilla/websocket"
	"io"
	"errors"
	"crypto/tls"
	"SocksOverWS/proxyconfig"
	"encoding/hex"
	"math/rand"
	"time"
	"strings"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const (
	socksAddr = "0.0.0.0:3000"
	testAddr  = "http://clients3.google.com/generate_204"
)

var socksListener net.Listener
var configuration proxyconfig.ProxyConfig
var TLSConfig tls.Config
var listening = false

func randomhost(t string) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	if t == "RANDOM" {
		domain := make([]byte, 6)
		rand.Read(domain)
		return hex.EncodeToString(domain) + ".cn"
	} else if t == "GFW" {
		return proxyconfig.GFWHosts[random.Intn(len(proxyconfig.GFWHosts)-1)]
	}
	return "" // shouldn't happen
}

func forward(connection net.Conn) {
	var dialer websocket.Dialer
	dialer = websocket.Dialer{
		TLSClientConfig: &TLSConfig,
	}
	if configuration.BypassType == "GFW" {
		dialer.TLSClientConfig.ServerName = randomhost(configuration.BypassType)
	} else if configuration.BypassType == "RANDOM" {
		dialer.TLSClientConfig.ServerName = randomhost(configuration.BypassType)
	}
	wsConnection, _, err := dialer.Dial(configuration.Addr, nil)
	if err != nil {
		log.Println("Connection Error", err.Error())
		return
	}
	go func() {
		for {
			if connection == nil {
				break
			}
			if !listening {
				break
			}
			_, message, err := wsConnection.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err) == true {
					break
				}
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) == true {
					break
				}
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				log.Println("Error while reading from websocket client", err.Error())
				break
			}
			connection.Write(message)
		}
		if connection != nil {
			connection.Close()
		}
		if wsConnection != nil {
			wsConnection.Close()
		}
		connection = nil
		wsConnection = nil
	}()
	go func() {
		for {
			if connection == nil {
				break
			}
			if !listening {
				break
			}
			buf := make([]byte, 8192)

			count, err := connection.Read(buf)
			if err != nil {
				if err == io.EOF && count > 0 {
					wsConnection.WriteMessage(websocket.BinaryMessage, buf[:count])
				}
				buf = nil
				break
			}
			if count > 0 {
				wsConnection.WriteMessage(websocket.BinaryMessage, buf[:count])
			}
			buf = nil
		}
		if connection != nil {
			connection.Close()
		}
		if wsConnection != nil {
			wsConnection.Close()
		}
		connection = nil
		wsConnection = nil
	}()
}
func getProxy(_ *http.Request) (*url.URL, error) {
	socksURL, err := url.Parse("socks5://" + socksAddr)
	return socksURL, err
}

func TestConnection() error {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: getProxy,
		},
		Timeout: 5 * time.Second,
	}
	req, err := httpClient.Get(testAddr)
	if err != nil {
		return err
	}
	req.Body.Close()
	return nil
}

func Run(config proxyconfig.ProxyConfig) error {
	TLSConfig.InsecureSkipVerify = !config.ValidateCert
	if config.EncryptionType == "aes128" {
		TLSConfig.CipherSuites = []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256}
	} else if config.EncryptionType == "chacha20" {
		TLSConfig.CipherSuites = []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305}
	}
	configuration = config
	socksServer, err := net.Listen("tcp", socksAddr)
	if err != nil {
		return errors.New("Failed to open socks port " + err.Error())
	}
	listening = true
	socksListener = socksServer
	var serverStarted sync.WaitGroup
	serverStarted.Add(1)
	go (func() {
		defer socksServer.Close()
		serverStarted.Done()
		for {
			if !listening {
				return
			}
			connection, _ := socksServer.Accept()
			go forward(connection)
		}
	})()
	serverStarted.Wait()
	return nil
}

func Stop() {
	listening = false
	socksListener.Close()
}
