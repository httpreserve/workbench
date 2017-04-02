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

var count int

func tformatOutput(ps processLog, response string) string {
	trStart := "<tr>"

	trFNAME := "<td>" + convertInterface(ps.lmap["response text"]) + "</td>"
	trVERSION := "<td>" + convertInterface(ps.lmap["analysis version text"]) + "</td>"
	trLINK := "<td>" + convertInterface(ps.lmap["link"]) + "</td>"
	trID := "<td>" + convertInterface(ps.lmap["archived"]) + "</td>"

	trEnd := "</tr>"

	response = response + trStart + trID + trFNAME + trVERSION + trLINK + trEnd
	return response 
}

func formatOutput(ps processLog, response string) string {
	count++
	//sfmt.Println(count)
	//<section id="section1">

	trStart := "<div><p>"

	trFNAME := convertInterface(ps.lmap["response text"]) + "<br/>"
	trVERSION := convertInterface(ps.lmap["analysis version text"]) + "<br/>"
	trLINK := convertInterface(ps.lmap["link"]) + "<br/>"
	trID := convertInterface(ps.lmap["archived"]) + "<br/>"

	trEnd := "</p></div>"

	response = response + trStart + trID + trFNAME + trVERSION + trLINK + trEnd

	fmt.Println(response)

	return response 
}

func aaformatOutput(ps processLog, response string) string {
	count++
	//sfmt.Println(count)
	//<section id="section1">

	trStart := "<div><pre>"

	trID := ps.js

	trEnd := "</pre></div>"

	response = response + trStart + trID + trEnd

	fmt.Println(response)

	return response 
}



func xformatOutput(ps processLog, response string) string {
	count++
	//sfmt.Println(count)
	//<section id="section1">

	trStart := fmt.Sprintf("<article><p>", count)

	trFNAME := convertInterface(ps.lmap["response text"]) + "<br/>"
	trVERSION := convertInterface(ps.lmap["analysis version text"]) + "<br/>"
	trLINK := convertInterface(ps.lmap["link"]) + "<br/>"
	trID := convertInterface(ps.lmap["archived"]) + "<br/>"

	trEnd := "</p></article>"

	response = response + trStart + trID + trFNAME + trVERSION + trLINK + trEnd
	return response 
}

func vformatOutput(ps processLog, response string) string {
	count++
	fmt.Println(count)
	//<section id="section1">

	trStart := fmt.Sprintf("<section id=\"section%d\"><p>", count)

	trFNAME := convertInterface(ps.lmap["response text"]) + "<br/>"
	trVERSION := convertInterface(ps.lmap["analysis version text"]) + "<br/>"
	trLINK := convertInterface(ps.lmap["link"]) + "<br/>"
	trID := convertInterface(ps.lmap["archived"]) + "<br/>"

	trEnd := "</p></section>"

	response = response + trStart + trID + trFNAME + trVERSION + trLINK + trEnd
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
