package httpmap

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

func CheckFiles(toCheck []string) (result bool, err error) {
	result = false

	files, err := Files()
	if err != nil {
		return
	}

	for _, check := range toCheck {
		if !slices.Contains(files, check) {
			return
		}
	}

	result = true
	return
}

func CheckKeys(file string, toCheck []string) (result bool, err error) {
	result = false

	keys, err := Keys(file)
	if err != nil {
		return
	}

	for _, check := range toCheck {
		if !slices.Contains(keys, check) {
			return
		}
	}

	result = true
	return
}

var Addr string = "http://127.0.0.1:8090"

func Files() (files []string, err error) {
	requestStr := Addr + "/files"

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

	files = strings.Split(string(data[1:len(data)-1]), " ")
	return
}

func Keys(file string) (keys []string, err error) {
	requestStr := Addr + fmt.Sprintf("/keys/%s", file)

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

	keys = strings.Split(string(data[1:len(data)-1]), " ")
	return
}

func Store(file string, key string, value string) (err error) {
	requestStr := Addr + fmt.Sprintf("/storage/%s/%s/%s", file, key, value)

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPut, requestStr, nil)
	if err != nil {
		return
	}

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

func Load(file string, key string) (value string, err error) {
	requestStr := Addr + fmt.Sprintf("/storage/%s/%s", file, key)

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

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	value = string(data)

	return
}
