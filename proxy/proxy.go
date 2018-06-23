package proxy

import (
	"net"
	"github.com/gorilla/websocket"
	"io"
	"errors"
	"crypto/tls"
	"log"
	"SocksOverWS/proxyconfig"
	"fmt"
)

var socksListener net.Listener
var configuration proxyconfig.ProxyConfig
var TLSConfig tls.Config
var listening = false

func randomhost(t string) string {
	return "hlol.wwww"
}

func forward(connection net.Conn) {
	var dialer websocket.Dialer
	dialer = websocket.Dialer{
		TLSClientConfig: &TLSConfig,
	}
	fmt.Println(dialer.TLSClientConfig)
	if configuration.BypassType == "GFW" {
		dialer.TLSClientConfig.ServerName = randomhost(configuration.BypassType)
	}
	wsConnection, _, err := dialer.Dial(configuration.Addr, nil)
	if err != nil {
		log.Println("Connection Error", err)
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
				log.Println(websocket.IsCloseError(err))
				log.Println("Error while reading from websocket client", err)
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

func Run(config proxyconfig.ProxyConfig) error {
	TLSConfig.InsecureSkipVerify = !config.ValidateCert
	if config.EncryptionType == "aes128" {
		TLSConfig.CipherSuites = []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256}
	} else if config.EncryptionType == "chacha20" {
		TLSConfig.CipherSuites = []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305}
	}
	configuration = config
	socksServer, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		return errors.New("Failed to open socks port " + err.Error())
	}
	listening = true
	socksListener = socksServer
	go (func() {
		defer socksServer.Close()
		for {
			if !listening {
				return
			}
			connection, _ := socksServer.Accept()
			go forward(connection)
		}
	})()
	return nil
}

func Stop() {
	listening = false
	socksListener.Close()
}
