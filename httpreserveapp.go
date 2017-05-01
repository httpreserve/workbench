package main

import (
	"flag"
	"fmt"
	"github.com/httpreserve/httpreserve"
	"log"
	"os"
	"time"
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
	csvout  bool
	webapp  bool

	//webapp config
	port string

	//list processing
	list string

	//throttle
	linkThrottle = 4
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

	//output method flags
	flag.BoolVar(&boltdb, "bolt", false, "Output to static BoltDB.")
	flag.BoolVar(&jsonout, "json", false, "Output to JSON.")
	flag.BoolVar(&csvout, "csv", false, "Output to CSV.")
	flag.BoolVar(&webapp, "webapp", false, "Output for analysis via webapp.")

	//other config parameters
	flag.StringVar(&port, "port", "", "Port to use for httpreserve webapp.")

	//create a list handler...
	flag.StringVar(&list, "list", "", "use test function while developing functionality.")
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

func getJSONFromLocal(link string, linklabel string) string {
	ls, err := httpreserve.GenerateLinkStats(link, linklabel, false)
	if err != nil {
		log.Println("Error retrieving linkstat JSON may be incorrect:", err)
	}
	js := httpreserve.MakeLinkStatsJSON(ls)

	// throttle requests to the server somehow...
	time.Sleep(500 * time.Millisecond)

	// return json...
	return js
}

func getLocalLink() {
	js := getJSONFromLocal(link, linklabel)
	fmt.Fprintln(os.Stderr, "Using httpreserve libs to retrieve data.")
	fmt.Fprintf(os.Stdout, "%s", js)
}

var htmcomplete bool
var starttime time.Time
var elapsedtime time.Duration

func programrunner() {

	if webapp {

		// processing time
		starttime = time.Now()

		// start web server and select{} below ensures
		// main doesn't complete...
		go webappRun()

		// next we need to get the data from the file...
		listHandler(webappHandler)

		// signal the server to start serving responses
		// todo: configure in parallel once this works
		htmcomplete = true

		htmpool[len(htmpool)-1].complete = true

		// don't return from function...
		select {}
	}

	if jsonout {
		fmt.Fprintf(os.Stdout, "%s", outputJSONHeader())
		listHandler(jsonHandler)
		outputjsonpool()
		fmt.Fprintf(os.Stdout, "%s", outputJSONFooter())
		return
	}

	if csvout {
		//output JSON header
		fmt.Fprintf(os.Stdout, "%s", outputCSVHeader())
		listHandler(csvHandler)
		return
	}

	if boltdb {
		openKVALBolt()
		defer closeKVALBolt()
		listHandler(boltdbHandler)
		return
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
		fmt.Fprintln(os.Stderr, "                                         [Optional -webapp] [Optional -port]")
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
