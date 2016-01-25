package fasthttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func JsonReqAndResHandler(url string, req interface{}, res interface{}, reqType string) (code int) {

	switch reqType {
	case "POST":
		var bodyBuffer *bytes.Buffer
		if req != nil {
			b, err := json.Marshal(&req)
			CheckError(err)
			log.Println("[FASTHTTP]the request is: ", string(b))
			bodyBuffer = bytes.NewBuffer(b)
		}
		response, err := http.Post(url, "application/json", bodyBuffer)
		return processHttpRes(response, err, &res)
	case "GET":
		response, err := http.Get(url)
		return processHttpRes(response, err, &res)
	default:
		log.Println("now in DEFAUTL swtich case...")
		return 0
	}
}

func processHttpRes(response *http.Response, err error, res interface{}) (code int) {
	log.Println("hey now in process")
	CheckError(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	CheckError(err)
	log.Printf("the status is: %s, and the body is: %s \n", response.Status, string(body))
	parseErr := json.Unmarshal(body, &res)
	CheckError(parseErr)
	return response.StatusCode
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln("fatal occur.", err.Error())
		panic(err)
	}
}
