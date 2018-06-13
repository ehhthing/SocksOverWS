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
var checksum []byte
func runGUIServer() string {
	static, err := fs.New()
	if err != nil {
		os.Exit(1)
	}
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
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
		available, sum, err := updater.Check()
		checksum = sum
		if err != nil {
			view.Eval(`alert("Error while checking for updates: ` + err.Error() + `")`)
		}
		if available {
			view.Eval(`window.receiveRPC({cmd: 'showUpdatePrompt'})`)
		}
	} else if cmd == "CONNECT" {
		if res["pac"].(string) == "GFW" {
			proxysettings.Set(server + "/gfw.pac")
		} else {
			proxysettings.Set(server + "/normal.pac")
		}
		proxy.Run(res["server"].(string), res["validateCertificate"].(string) == "true")
		view.Eval("window.receiveRPC({cmd: 'setConnectionStatus', status: true})")
	} else if cmd == "DISCONNECT" {
		proxysettings.Clear()
		proxy.Stop()
		view.Eval("window.receiveRPC({cmd: 'setConnectionStatus', status: false})")
	} else if cmd == "PAGECHANGE" {
		fmt.Println("Page Hijacked!")
		proxysettings.Clear()
		os.Exit(0)
	} else if cmd == "UPDATE" {
		view.Eval(`window.receiveRPC({cmd: 'showUpdateScreen'})`)
		go func() {
			err := updater.Update(checksum)
			if err != nil {
				view.Eval(`alert("Updater has failed, please download the latest version. \nError: ` + err.Error() + `")`)
				os.Exit(1)
			} else {
				view.Eval(`alert("Successfully updated, please start the proxy client again.")`)
				os.Exit(0)
			}
		}()
	}
}

func main() {
	var wait sync.WaitGroup
	wait.Add(1)
	server = runGUIServer()
	go (func() {
		view := webview.New(webview.Settings {
			Title:  "Socks over Websockets",
			URL:    server,
			Width:  300,
			Height: 225,
			ExternalInvokeCallback: handleRPC,
		})
		view.Run()
		proxysettings.Clear()
		wait.Done()
	})()
	wait.Wait()
}
