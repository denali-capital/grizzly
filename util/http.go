package util

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"net/http"
	"net/url"
)

// return errors with these functions?
// pointerize these functions?
func ParseUrlWithQuery(urlString string, values url.Values) string {
	url, err := url.Parse(urlString)
	if err != nil {
		log.Fatalln(err)
	}

	queryParams := url.Query()
	for param, valueArray := range values {
		for _, value := range valueArray {
			queryParams.Add(param, value)
		}
	}
	url.RawQuery = queryParams.Encode()

	return url.String()
}

func HttpGetAndGetBody(httpClient *http.Client, urlString string) map[string]interface{} {
	resp, err := httpClient.Get(urlString)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var bodyJson map[string]interface{}
	err = json.Unmarshal(body, &bodyJson)
	if err != nil {
		log.Fatalln(err)
	}

	return bodyJson
}

func DoHttpAndGetBody(httpClient *http.Client, request *http.Request) map[string]interface{} {
	resp, err := httpClient.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var bodyJson map[string]interface{}
	err = json.Unmarshal(body, &bodyJson)
	if err != nil {
		log.Fatalln(err)
	}

	return bodyJson
}