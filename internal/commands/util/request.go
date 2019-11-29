package util

import "net/http"

// URL url of the api
var URL = "https://api.testandset.com"

var apiVersion = "v1"

// APIGet run a GET request against the API
func APIGet(path string) (resp *http.Response, err error) {
	url := URL + "/" + apiVersion + "/" + path
	WriteVerboseMessage("Calling " + url)
	return http.Get(url)
}
