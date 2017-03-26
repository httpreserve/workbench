package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func makeHttpreservePOST() *http.Request {
	serviceURL := defaulthpserver
	resource := httpreserveFunction
	data := url.Values{}

	data.Set("url", link)
	data.Set("filename", linklabel)

	u, _ := url.ParseRequestURI(serviceURL)
	u.Path = resource

	r, _ := http.NewRequest("POST", u.String(), bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	return r
}

func makeRemoteRequest() {

	client := &http.Client{}
	resp, err := client.Do(makeHttpreservePOST())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error making client request")
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading body")
	}
	fmt.Fprintln(os.Stderr, "Using httpreserve web service to retrieve data.")
	fmt.Println(strings.Trim(string(responseData), " \n"))
}
