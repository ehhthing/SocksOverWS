package proxy

import (
	"net"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
)
var socks net.Listener
var http net.Listener
var listening = false
func forward(connection net.Conn, server string, connType string) {
	wsConnection, _, err := websocket.DefaultDialer.Dial(server + connType, nil)
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

func Run(server string) {
	listening = true
	socksServer, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		fmt.Println("Failed to create server", err)
	}

	httpServer, err := net.Listen("tcp", "localhost:3001")
	http = httpServer
	socks = socksServer
	if err != nil {
		fmt.Println("Failed to create server", err)
	}
	go (func() {
		defer socksServer.Close()
		for {
			if !listening {
				return
			}
			connection, _ := socksServer.Accept()
			go forward(connection, server, "socks")
		}
	})()
	go (func() {
		defer httpServer.Close()
		for {
			if !listening {
				return
			}
			connection, _ := httpServer.Accept()
			go forward(connection, server, "http")
		}
	})()
}

func Stop() {
	listening = false
	socks.Close()
	http.Close()
}