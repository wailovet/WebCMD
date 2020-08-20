package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	StartWeb("127.0.0.1:12302")
}

func StartWeb(addr string) {

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/cmd", func(writer http.ResponseWriter, request *http.Request) {
		cmd := exec.Command("powershell")
		h := websocket.Handler(func(ws *websocket.Conn) {

			cmd.Stdout = ws
			cmd.Stderr = ws
			cmd.Stdin = ws

			cmd.Run()

			log.Println("start")
			for {
				var reply string
				if err := websocket.Message.Receive(ws, &reply); err != nil {
					break
				}
			}
			cmd.Process.Kill()
		})
		h.ServeHTTP(writer, request)
	})

	_ = http.ListenAndServe(addr, nil)
}
