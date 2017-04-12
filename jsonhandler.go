package main

import (
	"fmt"
	"os"
)

func outputJSONHeader() string {
	var header string
	header = header + fmt.Sprintf("%s\n", "{")
	header = header + fmt.Sprintf("  \"%s\": \"%s\",\n", "title", "httpreserve")
	header = header + fmt.Sprintf("  \"%s\": \"%s\",\n", "description", "httpreserve client output")
	header = header + fmt.Sprintf("  \"%s\": %s\n", "data", "[")
	return header
}

func outputJSONFooter() string {
	var footer string
	footer = footer + fmt.Sprintf("%s\n%s", "]", "}")
	return footer
}

var jsonCount int

// TODO: consider more idiomatic approaches to achieving what we do here,
// that is, fmt.Println() is not really my approved approach (but it works (agile))
func jsonHandler(js string) {
	if js != "" {
		if jsonCount < linklen {
			fmt.Fprint(os.Stdout, js + ",")
		} else {
			fmt.Fprint(os.Stdout, js)
		}
	}
}
