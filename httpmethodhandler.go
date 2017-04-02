package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	//"net/url"
)

var complete bool
var indexlog	int
const fetchlen = 1		// select data from processedSlices in threes

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

func formatOutput(ps processLog, response string) string {

	trStart := "<div class=\"card\">"

	trLINK := "<b>" + convertInterface(ps.lmap["link"]) + "</b>"
	trRESP := "Response: " + convertInterface(ps.lmap["response code"]) + " " + convertInterface(ps.lmap["response text"])
	trSAVED := "Archived"  + convertInterface(ps.lmap["archived"])
	trFNAME := "Filename: " + convertInterface(ps.lmap["filename"])
	trSCREEN := "Screenshot: " + convertInterface(ps.lmap["screen shot"])
	trIAEARLIEST := "<b>IA Earliest:</b> " + convertInterface(ps.lmap["internet archive earliest"])
	trIALATEST := "<b>IA Latest:</b> " + convertInterface(ps.lmap["internet archive latest"])
	trIASAVE := "IA Save Link: " + convertInterface(ps.lmap["internet archive save link"])
	trIARESPCODE := "IA Response Code: " + convertInterface(ps.lmap["internet archive response code"])
	trIARESPONSETEXT := "IA Response Text: " + convertInterface(ps.lmap["internet archive response text"])

	trEnd := "</div>"
	trBR := "<br/>"

	response = response + trStart + trLINK + trBR + trRESP + trBR + trSAVED +
						trBR + trFNAME + trBR + trSCREEN + trBR + trIAEARLIEST + trBR + trIALATEST +
						trBR + trIASAVE + trBR + trIARESPCODE + trBR + trIARESPONSETEXT + trBR + trEnd

	//fmt.Println(ps.lmap)

	return response 
}

// Primary handler of all POST or GET requests to httpreserve
// pretty simple eh?!
func handleHttpreserve(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "%s", "time," + clock)
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
					indexlog = x+1
				}
				
				if complete {
					fmt.Fprintf(w, "false," + response)
				} else {
					fmt.Fprintf(w, "true," + response)
				}
			}
		} else {
			fmt.Println("no more data")
		}
	}
}
