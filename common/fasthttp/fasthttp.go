package fasthttp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func JsonReqAndResHandler(url string, req interface{}, res interface{}, reqType string) (code int) {
	var bodyBuffer io.Reader
	if req != nil {
		switch reqType {
		case "POST":
			bodyBuffer = getBufferReader(req)
		case "PUT":
			bodyBuffer = getBufferReader(req)
		default:
			// log.Println("now in DEFAUTL swtich case...")
		}
	}

	client := &http.Client{}
	log.Println(url)
	request, err := http.NewRequest(reqType, url, bodyBuffer)
	request.Header.Set("Content-Type", "application/json")
	CheckError(err)
	response, err := client.Do(request)
	return processHttpRes(response, err, &res)
}

func getBufferReader(req interface{}) *bytes.Buffer {
	if req != nil {
		b, err := json.Marshal(&req)
		CheckError(err)
		log.Println("[FASTHTTP]the request is: ", string(b))
		return bytes.NewBuffer(b)
	}
	return nil
}

func processHttpRes(response *http.Response, err error, res interface{}) (code int) {
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
