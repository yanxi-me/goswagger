package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/gobuffalo/packr"
)

func getAvailablePort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	return listener.Addr().(*net.TCPAddr).Port
}

func startServer(jsonPath string) {
	box := packr.NewBox("../web")
	http.Handle("/", http.FileServer(box))

	http.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		fh, err := os.Open(jsonPath)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(w, fh)
		if err != nil {
			log.Fatal(err)
		}
	})

	for port := 18000; ; port++ {
		timer := time.NewTimer(time.Second)
		go func() {
			<-timer.C
			url := "http://127.0.0.1:" + strconv.Itoa(port) + "/swagger.html"
			log.Println("URL: " + url)
			openBrowser(url)
		}()
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			timer.Stop()
			continue
		}
		break
	}
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: goswagger /swagger/json/path")
		return
	}

	jsonPath := os.Args[1]
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		fmt.Println(jsonPath, "not exist")
		return
	}

	startServer(jsonPath)
}
