package proxy

import (
	"net"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"errors"
	"crypto/tls"
)
var socksG net.Listener
var httpG net.Listener
var listening = false
func forward(connection net.Conn, server string, connType string, verifyCert bool) {
	var dialer websocket.Dialer
	if verifyCert == false {
		dialer = websocket.Dialer{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		dialer = websocket.Dialer{}
	}
	wsConnection, _, err := dialer.Dial(server + connType, nil)
	if err != nil {
		fmt.Println("Connection Error", err)
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
				if websocket.IsUnexpectedCloseError(err) == true {
					break
				}
				fmt.Println("Error while reading from websocket client", err)
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
			buf := make([]byte, 1024)

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

func Run(server string, verifyCert bool) error {
	if verifyCert == false {
		fmt.Println("Not verifying certificate")
	}
	listening = true
	socksServer, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		return errors.New("Failed to open socks port " + err.Error())
	}

	httpServer, err := net.Listen("tcp", "localhost:3001")
	httpG = httpServer
	socksG = socksServer
	if err != nil {
		return errors.New("Failed to open http proxy port " + err.Error())
	}
	go (func() {
		defer socksServer.Close()
		for {
			if !listening {
				return
			}
			connection, _ := socksServer.Accept()
			go forward(connection, server, "socks", verifyCert)
		}
	})()
	go (func() {
		defer httpServer.Close()
		for {
			if !listening {
				return
			}
			connection, _ := httpServer.Accept()
			go forward(connection, server, "http", verifyCert)
		}
	})()
	return nil
}

func Stop() {
	listening = false
	socksG.Close()
	httpG.Close()
}