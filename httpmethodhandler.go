package main

import (
	"fmt"
	"github.com/httpreserve/httpreserve"
	"github.com/httpreserve/wayback"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

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

const b64template = "{{ BASE64LOGO }}"
const screenshottemplate = "{{ SCREENSHOT CAPTION }}"

const column2 = `
<div class="column2">
	<figure class="screenshot">
		<img src="{{ BASE64LOGO }}" 
		width="250px" height="200px" alt="httpreserve"/></br>
		<figcaption><pre>screenshot for domain:</br>{{ SCREENSHOT CAPTION }}</figcaption>
	</figure>             
</div>`

const b64httpreservelogo = `
data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmc
vMjAwMC9zdmciIHdpZHRoPSI4IiBoZWlnaHQ9IjgiIHZpZXdCb3g9IjAgMCA4IDgiPg0KIC
A8cGF0aCBkPSJNMCAwdjFoOHYtMWgtOHptNCAybC0zIDNoMnYzaDJ2LTNoMmwtMy0zeiIgLz4NCjwvc3ZnPg==`

const placeholdercaption = "www.example.com"

const responseTable = `
	<table class="responsetable">
	<tr><td><b class="record">httpreserve record: </b></td><td class="two"><b><a target='_blank' class='httpreservelink' href='{{ DOMAIN }}'>{{ DOMAIN }}</td></tr>
	<tr><td>&nbsp;</td><td class="two">&nbsp;</td></tr>
	<tr><td>Response:</td><td class="two">{{ RESPONSE CODE }}</td></tr>
	<tr><td>Internet Archive:</td><td class="two">{{ ARCHVIED }}</td></tr>
	<tr><td>Filename:</td><td class="two">{{ FILENAME }}</td></tr>
	<tr><td>Title:</td><td class="two">{{ TITLE }}</td></tr>
	<tr><td>Content-type:</td><td class="two">{{ CONTENTTYPE }}</td></tr>		
	<tr><td>Wayback Earliest:</b></td><td class="two"><a target='_blank' class='httpreservelinkunder' href='{{ IA EARLY }}'>{{ IA EARLY HUMAN }}</a></td></tr>
	<tr><td>Wayback Latest:</b></td><td class="two"><a id='savelink{{ COUNT }}' target='_blank' class='httpreservelinkunder' href='{{ IA LATEST }}'>{{ IA LATEST HUMAN }}</a></td></tr>
	<tr><td>Wayback Save Link:</td><td class="two"><a target='_blank' class='httpreservelinkunder' href='javascript:saveToInternetArchive("{{ IA SAVE }}");'>{{ IA SAVE }}</a></td></tr>
	<tr><td>Wayback Response:</td><td class="two">{{ IA CODE }}</td></tr>
	<tr><td>Wayback Response Text:</td><td class="two">{{ IA TEXT }}</td></tr>
   </table>
`

const tbDomain = "{{ DOMAIN }}"
const tbCode = "{{ RESPONSE CODE }}"
const tbText = "{{ RESPONSE TEXT }}"
const tbArchived = "{{ ARCHVIED }}"
const tbFname = "{{ FILENAME }}"
const tbTitle = "{{ TITLE }}"
const tbContentType = "{{ CONTENTTYPE }}"
const tbSaveLink = "{{ IA SAVE }}"
const tbIACode = "{{ IA CODE }}"
const tbIAText = "{{ IA TEXT }}"

//dates
const tbIAEarly = "{{ IA EARLY }}"
const tbIALatest = "{{ IA LATEST }}"
const tbIAEarlyHuman = "{{ IA EARLY HUMAN }}"
const tbIALatestHuman = "{{ IA LATEST HUMAN }}"

var savecount int

const tbSaveID = "{{ COUNT }}"

func tableReplace(ps processLog) string {
	col1 := strings.Replace(responseTable, tbDomain, convertInterfaceHTML(ps.lmap["link"]), 2)
	col1 = strings.Replace(col1, tbCode, convertInterfaceHTML(ps.lmap["response code"]), 1)
	col1 = strings.Replace(col1, tbText, convertInterfaceHTML(ps.lmap["response text"]), 1)
	col1 = strings.Replace(col1, tbArchived, convertInterfaceHTML(ps.lmap["archived"]), 1)
	col1 = strings.Replace(col1, tbFname, convertInterfaceHTML(ps.lmap["filename"]), 1)
	col1 = strings.Replace(col1, tbTitle, convertInterfaceHTML(ps.lmap["title"]), 1)
	col1 = strings.Replace(col1, tbContentType, convertInterfaceHTML(ps.lmap["content-type"]), 1)
	col1 = strings.Replace(col1, tbDomain, convertInterfaceHTML(ps.lmap["screen shot"]), 1)

	earlydate := convertInterfaceHTML(ps.lmap["internet archive earliest"])
	latedate := convertInterfaceHTML(ps.lmap["internet archive latest"])
	col1 = strings.Replace(col1, tbIAEarly, earlydate, 1)
	col1 = strings.Replace(col1, tbIALatest, latedate, 1)
	col1 = strings.Replace(col1, tbIAEarlyHuman, wayback.GetHumanDate(earlydate), 1)
	col1 = strings.Replace(col1, tbIALatestHuman, wayback.GetHumanDate(latedate), 1)

	col1 = strings.Replace(col1, tbSaveLink, convertInterfaceHTML(ps.lmap["internet archive save link"]), 2)
	col1 = strings.Replace(col1, tbIACode, convertInterfaceHTML(ps.lmap["internet archive response code"]), 1)
	col1 = strings.Replace(col1, tbIAText, convertInterfaceHTML(ps.lmap["internet archive response text"]), 1)

	// make an id for the elements we output...
	id := fmt.Sprintf("%d", savecount)
	col1 = strings.Replace(col1, tbSaveID, id, 1)
	savecount++

	return col1
}

func addColumn1(columns string) string {
	return "<div class=\"column1\">" + columns + "</div>"
}

func addColumn2Default(ps processLog, columns string) string {
	snap := convertInterfaceHTML(ps.lmap["screen shot"])
	if snap == "" || strings.Contains(snap, httpreserve.SnapshotNotEnabled) ||
		strings.Contains(snap, httpreserve.GenerateSnapshotErr) {
		col2 := strings.Replace(column2, b64template, b64httpreservelogo, 1)
		col2 = strings.Replace(col2, screenshottemplate, convertInterfaceHTML(ps.lmap["link"]), 1)
		return columns + col2
	}

	// we may have a screenshot we can use...
	col2 := strings.Replace(column2, b64template, snap, 1)
	col2 = strings.Replace(col2, screenshottemplate, convertInterfaceHTML(ps.lmap["link"]), 1)

	return columns + col2
}

func addColumn2Live(columns string) string {
	return ""
}

func makeCardHTML(columns string) string {
	return "<div class=\"card\">" + columns + "</div>"
}

func formatOutput(ps processLog, response string) string {
	columns := tableReplace(ps)
	columns = addColumn1(columns)
	columns = addColumn2Default(ps, columns)
	return makeCardHTML(columns)
}

var pscopy []processLog
var outputcount int
var pscopyto int

var complete = false
var indexlog int

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
		buffersize := len(pscopy)

		// We want to maintain a whole copy of the list in memory to work
		// from, e.g. to update the indexes of. Do that here.
		if buffersize < processupdate && indexlog < processupdate {
			pscopyfrom := 0
			pscopy = pldatacopylen(&pscopyfrom, &pscopyto, processedSlices, 1)
		}

		//ensure neither buffer overruns the other...
		if buffersize > 0 && buffersize <= processupdate {

			if !complete && indexlog < processupdate {

				if pscopy[indexlog].complete == true {
					log.Println("received complete signal.")
					complete = true
				} else {
					response = formatOutput(pscopy[indexlog], response)
					log.Println(indexlog+1, "of", processupdate-1, "processed slices.") //one will always tell us to stop processing
				}
			}

			//finished processing what we've got, update indexlog
			//and only update indexlog if we've not got overunning buffers...
			indexlog++
		}

		// Let the client poll, unless a suitable exit condition is found...
		if complete {
			log.Println("Signalling client to stop polling.")
			fmt.Fprintf(w, "false•"+response)
		} else {
			fmt.Fprintf(w, "true•"+response)
		}
	}
}
