package main

import (
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"opensource/chaos/common"
	"opensource/chaos/domain"
	"opensource/chaos/domain/dao/mongo"
	"opensource/chaos/domain/dao/redis"
	oneClickDepRouter "opensource/chaos/modules/oneclickdep/router"
	webUtils "opensource/chaos/modules/oneclickdep/utils"
	"os"
	"path"
	"runtime"
	"time"
)

// templates := make(map[string]*template.Template) will occur "non-declaration statement outside function body" error
var templates = make(map[string]*template.Template)

var ARGS_STATIC_FILE_URL = flag.String("static", os.Getenv("GOPATH")+common.STATIC_DIR, "Static files address")
var ARGS_PROPERTIES_URL = flag.String("prop", os.Getenv("GOPATH")+common.PROP_FILE, "Properties files address")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("now begin to run server, please wait...")
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	route, _ := oneClickDepRouter.Route()
	api.SetApp(route)

	fileServer := http.FileServer(http.Dir(*ARGS_STATIC_FILE_URL))
	http.HandleFunc("/", entranceGuarder(indexHandler))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.ListenAndServe(":8080", nil)
	domain.Close()
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
	renderPage(w, "index", nil)
}

func renderPage(w http.ResponseWriter, tmpl string, values map[string]interface{}) {
	err := templates[tmpl+"."+common.TMPL_NAME].Execute(w, values)
	webUtils.CheckError(err)
}

func init() {
	flag.Parse()
	common.InitArgs(*ARGS_PROPERTIES_URL)
	mongo.MongoInit()
	redis.RedisInit()

	htmlPath := *ARGS_STATIC_FILE_URL + string(os.PathSeparator) + common.TMPL_NAME + string(os.PathSeparator)
	fileInfoArray, err := ioutil.ReadDir(htmlPath)
	webUtils.CheckError(err)
	var fileName, filePath string
	for _, fileInfo := range fileInfoArray {
		fileName = fileInfo.Name()
		if ext := path.Ext(fileName); ext != ("." + common.TMPL_NAME) {
			continue
		}
		filePath = htmlPath + fileName
		log.Println("loading template: " + filePath)
		t := template.Must(template.New(fileName).Delims("[[", "]]").ParseFiles(filePath))
		templates[fileName] = t
	}
}
