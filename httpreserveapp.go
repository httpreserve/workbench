package main

import (
	"flag"
	"fmt"
	"github.com/exponential-decay/httpreserve"
	"log"
	"net/http"
	"os"
)

var (
	demo      bool
	vers      bool
	demoport  string
	method    string
	link		 string
	linklabel string
)

const defaultPort = "2040"
const defaultMethod = http.MethodPost

func init() {
	flag.BoolVar(&demo, "demo", false, "Run demo server on port:2040 unless -port is set.")
	flag.BoolVar(&vers, "version", false, "Return httpreserve version.")
	flag.BoolVar(&vers, "v", false, "Return httpreserve version.")
	flag.StringVar(&demoport, "demoport", "", "Set a port to run httpreserve demo on localhost.")
	flag.StringVar(&method, "method", "", "Set a method to push queries through the demo, e.g. POST or GET.")
	flag.StringVar(&link, "link", "", "Seek the status of a single weblink.")
	flag.StringVar(&linklabel, "linklabel", "", "Annotate response with filename, or label.") 
}

func demosetup() {
	port := defaultPort
	meth := defaultMethod
	if demoport != "" {
		port = demoport
	}
	if method != "" {
		meth = method
	}
	err := httpreserve.DefaultServer(port, meth)
	if err != nil {
		log.Println("Error starting default server:", err)
	}
}

func getLink() {
	ls, err := httpreserve.GenerateLinkStats(link, "")
	if err != nil {
		log.Println("Error retrieving linkstat:", err)
		return
	}
	js := httpreserve.MakeLinkStatsJSON(ls)
	fmt.Fprintf(os.Stdout, "%s", js)
}

func programrunner() {
	if demo {
		//setup our demo server to communicate with
		demosetup()		
	}

	if link != "" {
		//retrieve data for a single link
		getLink()
	}
}

func main() {
	flag.Parse()
	if vers {
		fmt.Fprintf(os.Stderr, "%s\n", "httpreserve-app version information:")
		fmt.Fprintf(os.Stderr, "%s\n", httpreserve.VersionText())
		os.Exit(0)
	} else if flag.NFlag() <= 0 {
		fmt.Fprintln(os.Stderr, "Usage:  httpreserve-app [Optional -demo] [Optional -demoport] [Optional -method]")
		fmt.Fprintln(os.Stderr, "                        [Optional -link] [Optional -linklabel]")
		fmt.Fprintln(os.Stderr, "                        [Optional -version -v]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Output: [SERVER] 127.0.0.1:2040, [SERVER] 127.0.0.1:{port}")
		fmt.Fprintf(os.Stderr, "Output: [JSON] '%s ...'\n\n", "{ \"httpreserveanalysis\": \"x,y,z\" }")
		fmt.Fprintf(os.Stderr, "Output: [STRING] '%s ...'\n\n", httpreserve.VersionText())
		flag.Usage()
		os.Exit(0)
	}
	programrunner()
}
