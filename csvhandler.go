package main

import (
	"encoding/json"
	"fmt"
	"github.com/httpreserve/httpreserve"
	"os"
	"strings"
)

var csvHeader = []string{"id", "filename", "content-type", "title", "analysis version number", "analysis version text", "link",
	"response code", "response text", "screen shot", "internet archive latest", "internet archive earliest", "internet archive save link",
	"internet archive response code", "internet archive response text", "archived", "protocol error", "protocol error"}

func outputCSVHeader() string {
	var header string
	header = "\"" + strings.Join(csvHeader, "\",\"") + "\"" + "\n"
	return header
}

func outputCSVRow(lmap map[string]interface{}) string {

	var row []string
	for x := range csvHeader {
		if val, ok := lmap[csvHeader[x]]; ok {
			var v string
			switch val.(type) {
			case string:
				v = fmt.Sprintf("%s", val)
				v = strings.Replace(v, "\"", "'", -1)
				v = fmt.Sprintf("\"%s\"", v)
			case int:
				v = fmt.Sprintf("\"%d\"", val)
			case bool:
				v = fmt.Sprintf("\"%t\"", val)
			}
			row = append(row, v)
		} else {
			row = append(row, "\"\"")
		}
	}
	return strings.Join(row, ",")
}

// TODO: consider more idiomatic approaches to achieving what we do here,
// that is, fmt.Println() is not really my approved approach (but it works (agile))
func csvHandler(ch chan string) {

	var ls httpreserve.LinkStats

	ce := <-ch

	err := json.Unmarshal([]byte(ce), &ls)
	if err != nil {
		fmt.Fprintln(os.Stderr, "problem unmarshalling data.", err)
	}

	//fmt.Fprintf(os.Stdout, "%v\n", ls)

	// retrieve a map from the structure and write it out to the CSV
	lmap := storeStruct(ls, ce)
	if len(lmap) > 0 {
		fmt.Fprintf(os.Stdout, "%s\n", outputCSVRow(lmap))
	}
}
