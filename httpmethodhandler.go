package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

var complete bool
var indexlog int

const fetchlen = 1 // select data from processedSlices in threes

// For debug, we have this function here just in case we need
// to take a look at our request headers...
func prettyRequest(w http.ResponseWriter, r *http.Request) {
	req, _ := httputil.DumpRequest(r, false)
	fmt.Fprintf(w, "%s", req)
	return
}

// Min function for ints where Golang standard only handles
// int64...
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// make an id for the HTML elements we output...
var savecount int

// create a link that enables HTTPreserve to manage a
func makeSaveRequest(v interface{}) string {
	var val string
	id := fmt.Sprintf("%d", savecount)
	switch v.(type) {
	case string:
		if v != "" {
			val = fmt.Sprintf("%s", v)
			val = "<a class='httpreservelink' id=saveLink" + id + " target='_blank' href='javascript:saveToInternetArchive(\"" + val + "\");'>" + val + "</a>"
		}
	}
	savecount++
	return val
}

// convertInterface will help us pipe generic values from
// the deconstruction of httpreserve.LinkStats to a string for
// storage in BoltDB.
func convertInterfaceHTML(v interface{}) string {
	var val string
	switch v.(type) {
	case string:
		if v != "" {
			val = fmt.Sprintf("%s", v)
		} else {
			val = ""
		}
		if strings.Contains(val, "http") {
			val = "<a class='httpreservelink' target='_blank' href='" + val + "'>" + val + "</a>"
		}
	case int:
		val = fmt.Sprintf("%d", v)
	case bool:
		switch v {
		case true:
			val = "true"
		case false:
			val = "false"
		}
	}
	return val
}

const column2 = `
 <div class="column2">
	 <b>screenshot:</b>
	 <br/><br/>
    <img class="screenshot" src="https://github.com/exponential-decay/httpreserve/raw/master/src/images/httpreserve-logo.png"/>
 </div>`

func formatOutput(ps processLog, response string) string {

	trStart := "<div class=\"card\">"
	trColumn1 := "<div class=\"column1\">"

	trLINK := "<b class=\"record\">httpreserve record: </b><b>" + convertInterfaceHTML(ps.lmap["link"]) + "</b>"
	trRESP := "Response: " + convertInterfaceHTML(ps.lmap["response code"]) + " " + convertInterfaceHTML(ps.lmap["response text"])

	trSAVED := "Archived: " + "<span id='httpreserve-saved'>" + convertInterfaceHTML(ps.lmap["archived"]) + "</span>"

	trFNAME := "<b>Filename:</b> " + convertInterfaceHTML(ps.lmap["filename"])

	trCONTENTTYPE := "Content Type: " + convertInterfaceHTML(ps.lmap["content-type"])
	trTITLE := "Title: " + convertInterfaceHTML(ps.lmap["title"])

	/* Placeholder for screenshot output when the service works for us... */
	// trSCREEN := "Screenshot: " + convertInterfaceHTML(ps.lmap["screen shot"])

	trIAEARLIEST := "<b>IA Earliest:</b> " + convertInterfaceHTML(ps.lmap["internet archive earliest"])
	trIALATEST := "<b>IA Latest:</b> " + convertInterfaceHTML(ps.lmap["internet archive latest"])

	trIASAVE := "IA Save Link: " + makeSaveRequest(ps.lmap["internet archive save link"])

	trIARESPCODE := "IA Response Code: " + convertInterfaceHTML(ps.lmap["internet archive response code"])
	trIARESPONSETEXT := "IA Response Text: " + convertInterfaceHTML(ps.lmap["internet archive response text"])

	trColumn1End := "</div>"
	trEnd := "</div>"
	trBR := "<br/>"

	response = response + trStart + trColumn1 + trLINK + trBR + trBR + trRESP + trBR + trSAVED +
		trBR + trFNAME + trBR + trCONTENTTYPE + trBR + trTITLE + trBR + trIAEARLIEST + trBR + trIALATEST +
		trBR + trIASAVE + trBR + trIARESPCODE + trBR + trIARESPONSETEXT + trBR + trColumn1End + column2 + trEnd

	return response
}

var pscopy []processLog
var outputcount int
var pscopyto int

// Primary handler of all POST or GET requests to httpreserve
// pretty simple eh?!
func handleHttpreserve(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "%s", "time,"+clockOut())
		return
	case http.MethodPost:
		response := ""
		processupdate := len(processedSlices)

		// We want to maintain a whole copy of the list in memory to work
		// from, e.g. to update the indexes of. Do that here.
		if len(pscopy) < len(processedSlices) {
			pscopyfrom := 0
			//pscopy = pldatacopy(&pscopyfrom, &pscopyto, processedSlices)
			pscopy = pldatacopylen(&pscopyfrom, &pscopyto, processedSlices, 1)
		}

		if len(pscopy) > 0 {

			if !complete {
				limit := indexlog + (min(fetchlen, len(pscopy)))
				for x := indexlog; x < limit; x++ {
					if pscopy[x].complete == true {
						log.Println("received complete signal.")
						complete = true
						break
					}
					response = formatOutput(pscopy[x], response)
					indexlog = x + 1
					log.Println(indexlog, "of", processupdate, "processed slices")
				}

				if complete {
					fmt.Fprintf(w, "false•"+response)
				} else {
					fmt.Fprintf(w, "true•"+response)
				}
			}
		}
	}
}
