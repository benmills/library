package httpclient

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpRequest(method string, url string, data string) (int, string) {
	request, err := http.NewRequest(method, url, strings.NewReader(data))

	if err != nil {
		return 0, err.Error()
	}

	client := http.Client{}
	response, requestErr := client.Do(request)

	if requestErr != nil {
		return 0, requestErr.Error()
	}

	rawBody, bodyErr := ioutil.ReadAll(response.Body)

	if bodyErr != nil {
		panic(bodyErr)
	}

	return response.StatusCode, string(rawBody)
}

func Get(url string, data string) (int, string) {
	return HttpRequest("GET", url, data)
}

func Put(url string, data string) (int, string) {
	return HttpRequest("PUT", url, data)
}
