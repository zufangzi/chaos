package main

import (
	// "container/list"
	"cloud/server/entity"
	"encoding/json"
	"fmt"
	// "github.com/gorilla/mux"
	"github.com/tjz101/goprop"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

// Constant
const (
	STATIC_DIR = "/src/cloud/foreground"
	PROP_FILE  = "/src/cloud/resources/core.properties"
)

// global vars
var goCore GoCore

// templates := make(map[string]*template.Template) will occur "non-declaration statement outside function body" error
var templates = make(map[string]*template.Template)

// init config properties.
func init() {
	goCore = GoCore{""}
	prop := goprop.NewProp()
	prop.Read(os.Getenv("GOPATH") + PROP_FILE)
	goCore.MarathonAppsUrl = prop.Get("marathon.apps.url")

	htmlPath := os.Getenv("GOPATH") + STATIC_DIR + string(os.PathSeparator) + "html" + string(os.PathSeparator)
	fileInfoArray, err := ioutil.ReadDir(htmlPath)
	checkError(err)

	var fileName, filePath string
	for _, fileInfo := range fileInfoArray {
		fileName = fileInfo.Name()

		if ext := path.Ext(fileName); ext != ".html" {
			continue
		}
		filePath = htmlPath + fileName

		log.Println("loading template: " + filePath)

		t := template.Must(template.ParseFiles(filePath))
		// change delims to avoid the conflction with angularjs
		t.Delims("[", "]")
		templates[fileName] = t
	}

}

func main() {

	fileServer := http.FileServer(http.Dir(os.Getenv("GOPATH") + STATIC_DIR))
	http.Handle("/pic/", fileServer)
	http.Handle("/css/", fileServer)
	http.Handle("/font/", fileServer)
	http.Handle("/js/", fileServer)

	http.HandleFunc("/index", entranceGuarder(indexHandler))

	// get current-run apps' info
	http.HandleFunc("/apps", entranceGuarder(appsInfoHandler))

	http.ListenAndServe(":8080", nil)
}

func entranceGuarder(method http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now().UnixNano()
		defer func() {
			// func().(xx) means method return type cast
			// if in args type cast case. you can use string(xx) or xx.(string)
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				log.Println("catchable system error occur. " + e.Error())
			}

			log.Printf("the request: %s cost: %d ms\n", r.URL.RequestURI(), ((time.Now().UnixNano() - begin) / 1000000))
		}()
		method(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Welcome to DingDing Cloud, %s", html.EscapeString(r.URL.Path))
	// fmt.Fprintf(w, "\nyou host: %s", r.RemoteAddr)

	renderPage(w, "index", nil)

}

func appsInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("now enter in appinfo handler, remote addr: %s \n", html.EscapeString(r.RemoteAddr))
	resp, err := http.Get(goCore.MarathonAppsUrl)
	defer resp.Body.Close()
	checkError(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	// do not need to escape
	fmt.Fprintf(w, "Response is: \n %s", string(body))

	var appsInfo entity.MarathonAppsInfo

	json.Unmarshal(body, &appsInfo)

	fmt.Println(appsInfo)

	// break the
}

func renderPage(w http.ResponseWriter, tmpl string, values map[string]interface{}) {
	err := templates[tmpl+".html"].Execute(w, values)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatalln("fatal occur.", err.Error())
		panic(err)
	}
}

type GoCore struct {
	MarathonAppsUrl string
}
