package main

import (
	"net"
	"os"
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
	"SocksOverWS/proxyconfig"
	"log"
	"os/exec"
	"html"
)

var server string
var checksum []byte
var view webview.WebView
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
		err := proxy.Run(proxyconfig.ProxyConfig{
			Addr:           res["server"].(string),
			ValidateCert:   res["validateCertificate"].(string) == "true",
			EncryptionType: res["encryptionType"].(string),
			BypassType:     res["bypassType"].(string),
		})
		if err != nil {
			view.Eval("alert('Connection failed," + html.EscapeString(err.Error()) + "')")
			return
		}
		go (func() {
			log.Println("Testing connection")
			err = proxy.TestConnection()
			if err != nil {
				view.Dispatch(func() {
					view.Eval("alert('Failed to test connection," + html.EscapeString(err.Error()) + "')")
					view.Eval("window.receiveRPC({cmd: 'setConnectionStatus', status: false})")
				})
				return
			}
			view.Dispatch(func() {
				view.Eval("window.receiveRPC({cmd: 'setConnectionStatus', status: true})")
			})
		})()
	} else if cmd == "DISCONNECT" {
		proxysettings.Clear()
		proxy.Stop()
		view.Eval("window.receiveRPC({cmd: 'setConnectionStatus', status: false})")
	} else if cmd == "EXIT" {
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
	} else if cmd == "SHOW_LOG" {
		exec.Command("cmd","/C", "start", "powershell", "-Command", "Get-Content", "logs.log", "-Wait", "-Tail", "30").Run()
	}
}

func main() {
	var wait sync.WaitGroup
	wait.Add(1)
	server = runGUIServer()
	go (func() {
		view = webview.New(webview.Settings{
			Title:                  "Proxy Settings",
			URL:                    server,
			Width:                  290,
			Height:                 200,
			ExternalInvokeCallback: handleRPC,
		})
		view.Run()
		proxysettings.Clear()
		wait.Done()
	})()
	os.Remove("logs.log")
	logFile, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	wait.Wait()
}
