package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var linkLen int

// list handler to help us kick off some go channels
// we pass a first class function to help route our output
func listHandler(outputHandler func(ch chan string)) {
	ch := make(chan string)

	links, err := getList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file %s\n", err.Error())
	}

	linkLen = len(links)
	for l, f := range links {
		go channelLocalLink(l, f, ch)

		outputHandler(ch)

		//pause: TODO: Find a better pattern...
		time.Sleep(100 * time.Millisecond) //TODO: remove when throttling issues are solved
	}
}

func getList() (map[string]string, error) {
	var err error
	newlist := make(map[string]string)
	if list == "" {
		return linkmap, nil
	}
	newlist, err = readFile(list)
	return newlist, err
}

func readFile(l string) (map[string]string, error) {
	newlist := make(map[string]string)

	file, err := os.Open(list)
	if err != nil {
		return newlist, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), ",")
		if len(split) != 2 {
			fmt.Fprintf(os.Stderr, "ignoring: issue reading string from file: %s\n", scanner.Text())
		} else {
			newlist[strings.Trim(split[1], " ")] = strings.Trim(split[0], " ")
		}
	}

	if err := scanner.Err(); err != nil {
		return newlist, err
	}

	return newlist, nil
}

// demo linkmap...
var linkmap = map[string]string{
	"http://www.taupofest.co.nz/":         "nz govt",
	"http://www.siac.govt.nz":             "nz govt",
	"http://www.bbc.co.uk/news":           "bbc news",
	"http://www.bbc.co.uk/":               "bbc home",
	"http://www.bbc.co.uk/radio":          "bbc radio",
	"http://www.nationalarchives.gov.uk/": "tna",
	"http://www.google.com":               "",
	"http://google.com":                   "",
	"http://www.exponentialdecay.co.uk":   "",
	"http://www.archive.org":              "",
	"http://perma.cc":                     "",
	"http://wikipedia.org":                "",
}
