package httpmap

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var Addr string = "http://127.0.0.1:8090"

func Load(key string) (value []byte, err error) {
	requestStr := Addr + fmt.Sprintf("/storage/%s", key)

	var response *http.Response
	response, err = http.Get(requestStr)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("get '%s': %s", requestStr, http.StatusText(response.StatusCode))
		return
	}

	value, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	return
}

func Store(key string, value []byte) (err error) {
	requestStr := Addr + fmt.Sprintf("/storage/%s/%x", key, value)

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPut, requestStr, nil)
	if err != nil {
		return
	}

	// Fetch Request
	response, err := client.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("put '%s': %s", requestStr, http.StatusText(response.StatusCode))
		return
	}

	return
}

func File(fileName string) (err error) {
	requestStr := Addr + fmt.Sprintf("/file/%s", fileName)

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPut, requestStr, nil)
	if err != nil {
		return
	}

	// Fetch Request
	response, err := client.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("put '%s': %s", requestStr, http.StatusText(response.StatusCode))
		return
	}

	return
}

func Keys() (keys []string, err error) {
	requestStr := Addr + "/storage/keys"

	response, err := http.Get(requestStr)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("get '%s': %s", requestStr, http.StatusText(response.StatusCode))
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	keys = strings.Split(string(data), " ")

	return
}
