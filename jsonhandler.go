package main

import (
	"fmt"
	"os"
)

func channelLocalLink(link string, filename string, ch chan string) {
	ch <- libLink(link, filename)
}

//for now, for testing...
var linkmap = map[string]string{
	"http://www.bbc.co.uk/news":           "bbc news",
	"http://www.bbc.co.uk/":               "bbc home",
	"http://www.bbc.co.uk/radio":          "bbc radio",
	"http://www.nationalarchives.gov.uk/": "tna",
}

func outputHeader() string {
	var header string
	header = header + fmt.Sprintf("%s\n", "{")
	header = header + fmt.Sprintf("  \"%s\": \"%s\",\n", "title", "httpreserve")
	header = header + fmt.Sprintf("  \"%s\": \"%s\",\n", "description", "httpreserve client output")
	header = header + fmt.Sprintf("  \"%s\": %s\n", "data", "[")
	return header
}

func outputFooter() string {
	var footer string
	footer = footer + fmt.Sprintf("%s\n%s", "]", "}")
	return footer
}

// TODO: consider more idiomatic approaches to achieving what we do here,
// that is, fmt.Println() is not really my approved approach (but it works (agile))
func jsonHandler(ch chan string) {

	//output JSON header
	fmt.Fprintf(os.Stdout, "%s", outputHeader())

	//output JSON body
	var count int
	for range linkmap {
		count++
		ce := <-ch
		fmt.Print(ce)
		if count < len(linkmap) {
			fmt.Println(",")
		}
	}

	//output JSON footer
	fmt.Fprintf(os.Stdout, "%s", outputFooter())
}
