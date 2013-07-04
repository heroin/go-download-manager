package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"text/template"
	"io"
	"bufio"
	"io/ioutil"
	"strings"
)

const (
	DOWNLOAD_PATH = "/www/domain/download.heroin.so"
	SERVER        = "GO-DOWNLOAD-MANAGER"
	PORT          = 12323
	CMD_LINE      = "wget"
)

type Context struct {
	data interface{}
}

type File struct {
	Name string
	Date string
	Dir  bool
}

func (context *Context) view(write io.Writer, html string) {
	cursor := template.Must(template.ParseFiles(fmt.Sprintf("views/%s.html", html)))
	cursor.Execute(write, context.data)
}

func index(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	requestPath := request.URL.Path
	if requestPath != "/" {
		n := len(requestPath)
		http.ServeFile(out, request, fmt.Sprintf("static/%s", requestPath[1:n]))
	} else {
		app := Context {}
		app.view(out, "index")
	}
}

func list(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	request.ParseForm()
	path := request.FormValue("path")
	var data []File
	if strings.HasPrefix(path, "/") && strings.Index(path, "../") < 0 {
		files, _ := ioutil.ReadDir(fmt.Sprintf("%s%s", DOWNLOAD_PATH, path))
		for _, i := range files {
			data = append(data, File {
					Name: i.Name(),
					Date: i.ModTime().Format("2006-01-02 15:04:05"),
					Dir: i.IsDir(),
				})
		}
	}
	cursor := template.Must(template.New("text").Parse(`<?xml version="1.0" encoding="UTF-8"?><root>{{if .}}{{range .}}<file><name>{{.Name}}</name><date>{{.Date}}</date><dir>{{.Dir}}</dir></file>{{end}}{{end}}</root>`))
	cursor.Execute(out, data)
}

func download(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	request.ParseForm()
	url := request.FormValue("url")
	name := request.FormValue("name")
	path := request.FormValue("path")

	if strings.TrimSpace(url) != "" {
		if strings.TrimSpace(name) != "" {
			if strings.Index(name, "../") < 0 && strings.Index(name, "/") < 0 {
				go Download(url, generatePath(path), name)
				fmt.Fprintf(out, "{\"result\":\"Success\", \"code\":1}")
			} else {
				go Download(url, generatePath(path), "")
				fmt.Fprintf(out, "{\"result\":\"Success\", \"code\":1}")
			}
		} else {
			go Download(url, generatePath(path), "")
			fmt.Fprintf(out, "{\"result\":\"Success\", \"code\":1}")
		}
	} else {
		fmt.Fprintf(out, "{\"result\":\"Error\", \"code\":-1}")
	}
}

func batchDownload(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	request.ParseForm()
	urls := request.FormValue("urls")
	path := request.FormValue("path")
	reader := bufio.NewReader(strings.NewReader(urls))
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}else {
			go Download(string(line), generatePath(path), "")
		}
	}
	fmt.Fprintf(out, "{\"result\":\"Success\", \"code\":1}")
}

func remove(out http.ResponseWriter, request *http.Request) {
	out.Header().Set("Server", SERVER)
	request.ParseForm()
	fmt.Println(request.Form["file"])
	fmt.Println(request.FormValue("file"))
	fmt.Println(request.Form["dir"])
	fmt.Println(ioutil.ReadDir("./"))
	files, _ := ioutil.ReadDir("./")
	for i := range files {
		fmt.Println(files[i].Name())
	}
	fmt.Fprintf(out, "rm\r\n")
}

func Download(url string, path string, name string) {
	runtime.Gosched()
	log.Printf("download start %s \n", url)
	cmd := exec.Command(CMD_LINE)
	if strings.TrimSpace(name) != "" {
		cmd = exec.Command(CMD_LINE, "-O", fmt.Sprintf("%s%s", path, name), url)
	} else {
		cmd = exec.Command(CMD_LINE, "-P", path, url)
	}
	err := cmd.Run()
	if err != nil {
		log.Printf("download [error] path=%s, ", url, err)
	}
	runtime.GC()
	log.Printf("download over %s \n", url)
}

func generatePath(path string) string {
	if strings.TrimSpace(path) != "" {
		if !strings.HasPrefix(path, "/") && strings.Index(path, "../") < 0 {
			return fmt.Sprintf("%s/%s/", DOWNLOAD_PATH, path)
		}
		return DOWNLOAD_PATH
	}
	return DOWNLOAD_PATH
}

func main() {
	runtime.GOMAXPROCS(5)
	http.HandleFunc("/", index)
	http.HandleFunc("/list", list)
	http.HandleFunc("/rm", remove)
	http.HandleFunc("/download", download)
	http.HandleFunc("/batch/download", batchDownload)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
