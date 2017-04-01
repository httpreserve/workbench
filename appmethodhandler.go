package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	//"net/url"
)

const requestedURL = "url"
var count = 10
var no = 0
var tout = true

// For debug, we have this function here just in case we need
// to take a look at our request headers...
func prettyRequest(w http.ResponseWriter, r *http.Request) {
	req, _ := httputil.DumpRequest(r, false)
	fmt.Fprintf(w, "%s", req)
	return
}

// Primary handler of all POST or GET requests to httpreserve
// pretty simple eh?!
func handleHttpreserve(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//["Saab", "Volvo", "BMW"];
		fmt.Fprintf(w, "%s", "time," + clock)
		return
	case http.MethodPost:
		switch tout {
		case false:
			fmt.Fprintf(w, "%s", "false" + "," + clock)
		case true:
			fmt.Fprintf(w, "%s", "true" + "," + clock)
		}

		no++
		if no >= count {
			tout = false
		}
		return 
	}
}
