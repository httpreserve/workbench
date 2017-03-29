package main

import (
	"fmt"
	"os"
)

func channelLocalLink(link string, filename string, ch chan string) {
	ch <- libLink(link, filename)
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
