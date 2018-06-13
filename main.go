package main

import (
	"net"
	"os"
	"fmt"
	"github.com/zserge/webview"
	"github.com/rakyll/statik/fs"
	"net/http"
	_ "./statik"
	"strconv"
	"./proxysettings"
	"encoding/json"
	"sync"
	"SocksOverWS/proxy"
	"SocksOverWS/updater"
)
var server string

func runGUIServer() string {
	static, err := fs.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	http.Handle("/", http.FileServer(static))
	go (func() {
		http.Serve(listener, nil)
	})()
	return "http://127.0.0.1:" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
}

func handleRPC(view webview.WebView, data string) {
	var res map[string]interface{}
	json.Unmarshal([]byte(data), &res)
	cmd := res["action"].(string)
	if cmd == "READY" {
		fmt.Println("GUI Ready")
	} else if cmd == "CONNECT" {
		if res["pac"].(string) == "GFW" {
			proxysettings.Set(server + "/gfw.pac")
		} else {
			proxysettings.Set(server + "/normal.pac")
		}
		proxy.Run(res["server"].(string))
		view.Eval("window.receiveRPC({cmd: 'setConnectionStatus', status: true})")
	} else if cmd == "DISCONNECT" {
		proxysettings.Clear()
		proxy.Stop()
		view.Eval("window.receiveRPC({cmd: 'setConnectionStatus', status: false})")
	} else if cmd == "PAGECHANGE" {
		fmt.Println("Page Hijacked!")
		proxysettings.Clear()
		os.Exit(0)
	}
}
func main() {
	updater.Check()
	return
	var wait sync.WaitGroup
	wait.Add(1)
	server = runGUIServer()
	fmt.Println(server)
	go (func() {
		view := webview.New(webview.Settings{
			Title:  "Socks over Websockets",
			URL:    server,
			Width:  300,
			Height: 215,
			ExternalInvokeCallback: handleRPC,
		})
		view.Run()
		proxysettings.Clear()
		wait.Done()
	})()
	wait.Wait()
}
