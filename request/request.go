package request

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"gitlab.com/marioskill/configuration"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func Do(query interface{}) (Result, error) {
	var http_response Result
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(query)
	resp, err := http.Post(configuration.Mesos_Server_url_API, "application/json", b)

	if err != nil {
		http_response.Http_code = 503
		http_response.Http_msg = "Mesos is not working"

	} else {
		codehttp := strings.Split(resp.Status, " ")
		http_response.Http_code, _ = strconv.Atoi(codehttp[0])
		http_response.Http_msg = codehttp[1]

		body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))
		http_response.Response = body

	}

	return http_response, err
}
