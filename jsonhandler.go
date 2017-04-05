package main

import (
	"fmt"
	"os"
)

func channelLocalLink(link string, filename string, ch chan string) {
	ch <- libLink(link, filename)
}

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
func jsonHandler(ce string) {

	//output JSON header
	//fmt.Fprintf(os.Stdout, "%s", outputHeader())

	//output JSON body
	//var count int
	//for range linkmap {
		//count++
		//ce := <-ch
	
	jsonCount++
	if ce != "" {
		if jsonCount < linkLen {
			fmt.Print(ce + ",")
		} else {
			fmt.Print(ce)
			fmt.Fprintf(os.Stderr, "no comma")
		}
	}

	fmt.Fprintf(os.Stderr, "%d, %d\n", jsonCount, linkLen)

	//if count < len(linkmap) {
	//mt.Println(",")
	//}
	//}

	//output JSON footer
	//fmt.Fprintf(os.Stdout, "%s", outputFooter())
}
