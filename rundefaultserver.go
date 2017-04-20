package main

import (
	"fmt"
	"github.com/httpreserve/httpreserve"
	"github.com/httpreserve/simplerequest"	
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func launchRemoteBackgroundProcess() {
	cmd := exec.Command("httpreserve-app", "-demo")
	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting myself in the background")
	}
	time.Sleep(2)
}

func testDefaultServer() bool {
	// no cost to us to try and start ourselves...
	launchRemoteBackgroundProcess()

	// setup a simple call to the server to look for a 200 OK response
	sr, err := simplerequest.Create(http.MethodHead, defaulthpserver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return false
	}

	// send a request to the server to see if it's there...
	ls, err := httpreserve.HTTPFromSimpleRequest(sr)
	if err != nil {
		if strings.Contains(err.Error(), "getsockopt: connection refused") {
			fmt.Fprintln(os.Stderr, "httpreserve server is not running.")
		}
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return false
	}

	// if our response is anything other than 200 OK return to exit...
	if ls.ResponseCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "httpreserve server is not running.")
		return false
	}

	return true
}
