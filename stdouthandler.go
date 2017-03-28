package main

import (
	"fmt"
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

// TODO: consider more idiomatic approaches to achieving what we do here,
// that is, fmt.Println() is not really my approved approach (but it works (agile))
func toStdout(ch chan string) {
	fmt.Println("{")
	fmt.Println("\"title\": \"httpreserve client example\",")
	fmt.Println("\"data\": [")

	var count int
	for range linkmap {
		count += 1
		ce := <-ch
		fmt.Print(ce)
		if count < len(linkmap) {
			fmt.Println(",")
		}
	}
	fmt.Println("]\n}")
}

func jsonhandler() {
	ch := make(chan string)
	for l, f := range linkmap {
		go channelLocalLink(l, f, ch)
	}
	toStdout(ch)
}
