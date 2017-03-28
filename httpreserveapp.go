package main

import (
	"flag"
	"fmt"
	"github.com/exponential-decay/httpreserve"
	"log"
	"os"
)

var (
	demo       bool
	vers       bool
	demoport   string
	demomethod string
	link       string
	linklabel  string
	remote     bool

	//output methods
	boltdb  bool
	jsonout bool
	webapp  bool

	//testing
	test bool
)

func init() {
	flag.BoolVar(&demo, "demo", false, "Run demo server on port:2040 unless -port is set.")
	flag.BoolVar(&vers, "version", false, "Return httpreserve version.")
	flag.BoolVar(&vers, "v", false, "Return httpreserve version.")

	//flags to return a single result
	flag.StringVar(&link, "link", "", "Seek the status of a single weblink.")
	flag.StringVar(&linklabel, "linklabel", "", "Annotate response with filename, or label.")

	//demo configuration
	flag.StringVar(&demoport, "demoport", "", "Set a port to run httpreserve demo on localhost.")
	flag.StringVar(&demomethod, "demomethod", "", "Set a method to push queries through the demo, e.g. POST or GET.")

	//retireve stats from web service
	flag.BoolVar(&remote, "remote", false, "Send requests to remote connection.")

	//testing flags for testing whatever we're working on at present
	flag.BoolVar(&test, "test", false, "use test function while developing functionality.")

	//output method flags
	flag.BoolVar(&boltdb, "bolt", false, "Output to static BoltDB.")
	flag.BoolVar(&jsonout, "json", false, "Output to JSON.")
	flag.BoolVar(&webapp, "webapp", false, "Output for analysis via webapp.")
}

func demosetup() {
	port := defaultPort
	method := defaultMethod
	if demoport != "" {
		port = demoport
	}
	if demomethod != "" {
		method = demomethod
	}
	err := httpreserve.DefaultServer(port, method)
	if err != nil {
		log.Println("Error starting default server:", err)
		os.Exit(1)
	}
}

func getRemoteLink() {
	var remotelinkExists bool
	if requestType == remoteRequest {
		remotelinkExists = testDefaultServer()
	}
	if remotelinkExists {
		makeRemoteRequest()
	}
}

func libLink(link string, linklabel string) string {
	ls, err := httpreserve.GenerateLinkStats(link, linklabel)
	if err != nil {
		log.Println("Error retrieving linkstat:", err)
		return ""
	}
	js := httpreserve.MakeLinkStatsJSON(ls)
	return js
}

func getLocalLink() {
	js := libLink(link, linklabel)
	fmt.Fprintln(os.Stderr, "Using httpreserve libs to retrieve data.")
	fmt.Fprintf(os.Stdout, "%s", js)
}

func programrunner() {

	if jsonout {
		jsonhandler()
		return
	}

	if boltdb {
		boltdbHandler()
		return
	}

	if webapp {
		fmt.Fprintln(os.Stderr, "Handling webapp here.")
	}

	if demo {
		//setup our demo server to communicate with
		demosetup()
	}

	if link != "" && remote {
		//retrieve data for a single link
		getRemoteLink()
	}

	if link != "" && !remote {
		getLocalLink()
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
		fmt.Fprintln(os.Stderr, "                        [Optional -link] [Optional -linklabel] [Optional -remote]")
		fmt.Fprintln(os.Stderr, "                        [Optional -list] [Optional -json]")
		fmt.Fprintln(os.Stderr, "                                         [Optional -bolt]")
		fmt.Fprintln(os.Stderr, "                                         [Optional -webapp]")
		fmt.Fprintln(os.Stderr, "                        [Optional -version -v]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Output: [SERVER] 127.0.0.1:2040, [SERVER] 127.0.0.1:{port}")
		fmt.Fprintf(os.Stderr, "Output: [JSON] '%s ...'\n", "{ \"httpreserveanalysis\": \"x,y,z\" }")
		fmt.Fprintf(os.Stderr, "Output: [VERSION] '%s ...'\n", httpreserve.VersionText())
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(0)
	}
	programrunner()
}
