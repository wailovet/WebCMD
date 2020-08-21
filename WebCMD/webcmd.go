package WebCMD

import (
	"flag"
	"golang.org/x/net/websocket"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"net/http"
	"os/exec"
)

func StartWeb(addr string) {
	var dev bool
	flag.BoolVar(&dev, "dev", false, "")
	flag.Parse()
	if dev {
		http.Handle("/", http.FileServer(http.Dir("static")))
	} else {
		http.Handle("/", http.FileServer(assetFS()))
	}

	http.HandleFunc("/cmd", func(writer http.ResponseWriter, request *http.Request) {
		cmds := request.URL.Query().Get("s")
		dir := request.URL.Query().Get("dir")
		h := websocket.Handler(func(ws *websocket.Conn) {

			cmd := exec.Command("cmd", "/c", cmds)
			cmd.Dir = dir
			gbks := transform.NewWriter(ws, simplifiedchinese.GBK.NewDecoder())
			cmd.Stdout = gbks
			cmd.Stderr = gbks

			cmd.Run()
		})
		h.ServeHTTP(writer, request)
	})

	_ = http.ListenAndServe(addr, nil)
}
