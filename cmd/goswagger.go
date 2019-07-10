package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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

func startServer(jsonPath string, port int) {
	box := packr.NewBox("../web")
	jsonName := filepath.Base(jsonPath)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			html, _ := box.FindString("swagger.html")
			fmt.Fprint(w, strings.ReplaceAll(html, "__swaggerUrl__", "/"+jsonName))
		} else {
			http.FileServer(box).ServeHTTP(w, req)
		}
	})

	http.HandleFunc("/"+jsonName, func(w http.ResponseWriter, req *http.Request) {
		fh, err := os.Open(jsonPath)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(w, fh)
		if err != nil {
			log.Fatal(err)
		}
	})

	for ; ; port++ {
		timer := time.NewTimer(time.Second)
		go func() {
			<-timer.C
			url := "http://127.0.0.1:" + strconv.Itoa(port) + "/"
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
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: goswagger your/swagger/json/path [port]")
		return
	}

	jsonPath := os.Args[1]
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		fmt.Println(jsonPath, "not exist")
		return
	}

	port := 18000 // default port
	if len(os.Args) == 3 {
		var err error
		if port, err = strconv.Atoi(os.Args[2]); err != nil {
			fmt.Println("port must be a number")
			return
		}
	}

	startServer(jsonPath, port)
}
