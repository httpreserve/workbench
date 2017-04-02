package main

import (
	"fmt"
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
			val = "<a class='httpreservelink' href='" + val + "'>" + val + "</a>"
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
	trSAVED := "Archived: " + convertInterfaceHTML(ps.lmap["archived"])
	trFNAME := "Filename: " + convertInterfaceHTML(ps.lmap["filename"])
	trSCREEN := "Screenshot: " + convertInterfaceHTML(ps.lmap["screen shot"])
	trIAEARLIEST := "<b>IA Earliest:</b> " + convertInterfaceHTML(ps.lmap["internet archive earliest"])
	trIALATEST := "<b>IA Latest:</b> " + convertInterfaceHTML(ps.lmap["internet archive latest"])
	trIASAVE := "IA Save Link: " + convertInterfaceHTML(ps.lmap["internet archive save link"])
	trIARESPCODE := "IA Response Code: " + convertInterfaceHTML(ps.lmap["internet archive response code"])
	trIARESPONSETEXT := "IA Response Text: " + convertInterfaceHTML(ps.lmap["internet archive response text"])

	trColumn1End := "</div>"
	trEnd := "</div>"
	trBR := "<br/>"

	response = response + trStart + trColumn1 + trLINK + trBR + trBR + trRESP + trBR + trSAVED +
		trBR + trFNAME + trBR + trSCREEN + trBR + trIAEARLIEST + trBR + trIALATEST +
		trBR + trIASAVE + trBR + trIARESPCODE + trBR + trIARESPONSETEXT + trBR + trColumn1End + column2 + trEnd

	return response
}

// Primary handler of all POST or GET requests to httpreserve
// pretty simple eh?!
func handleHttpreserve(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "%s", "time,"+clock)
		return
	case http.MethodPost:
		response := ""
		if len(processedSlices) > 0 {
			if !complete {
				limit := indexlog + (min(fetchlen, len(processedSlices)))
				for x := indexlog; x < limit; x++ {
					fmt.Println(x)
					if processedSlices[x].complete == true {
						complete = true
						break
					}
					response = formatOutput(processedSlices[x], response)
					indexlog = x + 1
				}

				if complete {
					fmt.Fprintf(w, "false,"+response)
				} else {
					fmt.Fprintf(w, "true,"+response)
				}
			}
		}
	}
}
