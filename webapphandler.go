package main

import (
	"encoding/json"
	"github.com/httpreserve/httpreserve"
	"log"
	"sync"
	"time"
)

// This structure is used to communicate with the server
// we may also use some static storage in the form of Bolt DB
// the final signal to the webapp will be a empty payload
// but with the complete flag set to true so that we know there
// is no more work to be processed. ls contains a link stat
// data structure if we can recreate one from the JSON we receive
// else the js variable will contain a single JSON document to
// be processed.
type processLog struct {
	complete bool
	ls       httpreserve.LinkStats
	js       string
	lmap     map[string]interface{}
}

func clockOut() string {
	t := time.Now()
	return t.Format("Mon Jan _2 15:04:05 2006")
}

// webapprun lets us start the server for the user to access
func webappRun() {
	//defer serverWG.Done()

	// pause to let our buffers begin to fill...
	// TODO: look for safer, more idiomatic ways to solve...
	//time.Sleep(3 * time.Second)
	log.Println("Server starting on localhost: http://127.0.0.1:2041")

	if port == "" {
		port = "2041"
	}

	err := DefaultServer(port)
	if err != nil {
		log.Println(err)
	}

}

var lpcopyfrom, lpcopyto int
var processedSlices []processLog

func processlinkpool(wg *sync.WaitGroup) {
	defer wg.Done()

	res := tsdatacopy(&lpcopyfrom, &lpcopyto, linkpool)

	if len(res) > 0 {

		var ls httpreserve.LinkStats

		for x := range res {
			ce := res[x]
			err := json.Unmarshal([]byte(ce), &ls)
			if err != nil {
				log.Println("Problem unmarshalling data.", err)
			}

			// retrieve a map from the structure and write it out to the
			// http server...
			lmap := storeStruct(ls, ce)
			if len(lmap) > 0 {
				var ps processLog
				ps.js = ce
				ps.ls = ls
				ps.lmap = lmap
				processedSlices = append(processedSlices, ps)
			}
		}

		if pscomplete {
			var ps processLog
			ps.complete = true
			processedSlices = append(processedSlices, ps)
		}
	}
}

var linkpool []string

func makelinkpool(js string, wg *sync.WaitGroup) {
	defer wg.Done()
	linkpool = append(linkpool, js)
}

// temporary webappHandler that gets the app working reliably...
func webappHandler(js string) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go makehtmpool(js, wg)
	wg.Wait()
	return
}

var htmpool []processLog

func makehtmpool(js string, wg *sync.WaitGroup) {
	defer wg.Done()

	var ls httpreserve.LinkStats
	err := json.Unmarshal([]byte(js), &ls)
	if err != nil {
		log.Println("Problem unmarshalling data.", err)
		return
	}

	// retrieve a map from the structure and write it out to the
	// http server...
	lmap := storeStruct(ls, js)
	if len(lmap) > 0 {
		var ps processLog
		ps.js = js
		ps.ls = ls
		ps.lmap = lmap
		htmpool = append(htmpool, ps)
	}

	log.Println("processed record:", len(htmpool))

	return
}

// webappHanlder enables us to establish the web server and create
// the structures we need to present our data to the user...
func concurrentversionNotWorkingWebappHandler(js string) {

	wg := new(sync.WaitGroup)

	//TODO: Understand the essence of wait groups where the ordering
	//of these, slightly backwards works... does processlinkpool just
	//wait for makelinkpool????
	wg.Add(1)
	go processlinkpool(wg)

	wg.Add(1)
	go makelinkpool(js, wg)

	wg.Wait()
	return
}
