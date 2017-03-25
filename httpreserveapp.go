package main

import (
	"flag"
	"fmt"
	"github.com/exponential-decay/httpreserve"
	"net/http"
	"os"
)

var (
	demo     bool
	vers     bool
	demoport string
	method   string
)

const defaultPort = "2040"
const defaultMethod = http.MethodPost

func init() {
	flag.BoolVar(&demo, "demo", false, "Run demo server on port:2040 unless -port is set.")
	flag.BoolVar(&vers, "version", false, "Return httpreserve version.")
	flag.StringVar(&demoport, "demoport", "", "Set a port to run httpreserve demo on localhost.")
	flag.StringVar(&method, "method", "", "Set a method to push queries through the demo, e.g. POST or GET.")
}

func programrunner() {
	if demo {
		port := defaultPort
		meth := defaultMethod
		if demoport != "" {
			port = demoport
		}
		if method != "" {
			meth = method
		}
		httpreserve.DefaultServer(port, meth)
	}
}

func main() {
	flag.Parse()
	var verstring = "version"
	if vers {
		fmt.Fprintf(os.Stderr, "%s %s\n", verstring, httpreserve.VersionText())
		os.Exit(0)
	} else if flag.NFlag() <= 0 {
		fmt.Fprintln(os.Stderr, "Usage:  httpreserve-app [Optional -demo] [Optional -demoport] [Optional -method]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Output: [SERVER] 127.0.0.1:2040, [SERVER] 127.0.0.1:{port}")
		fmt.Fprintf(os.Stderr, "Output: [STRING] '%s ...'\n\n", httpreserve.VersionText())
		flag.Usage()
		os.Exit(0)
	}
	programrunner()
}
