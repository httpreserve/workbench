package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

var linklen int

// list handler to help us kick off some go channels
// we pass a first class function to help route our output
func listHandler(outputHandler func(js string)) {

	//fmt.Println("in list handler")

	link := make(chan map[string]string)
	results := make(chan string)

	wg := new(sync.WaitGroup)

	// 10 chunks of work..?
	for w := 0; w <= 20; w++ {
		wg.Add(1)
		go getJSON(link, results, wg)
	} 

	// Create a link map for output to the output handlers
	go func() {
		if list == "" {
			// Use demo list...
			for k, v := range linkmap {
				l := make(map[string]string)
				l[k] = v
				link <- l
			}
		} else {
			// Read all the lines in a file...
			file, err := os.Open(list)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error with scanner: %s\n", err.Error())
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				l := make(map[string]string)
				split := strings.Split(scanner.Text(), ",")
				if len(split) != 2 {
					fmt.Fprintf(os.Stderr, "ignoring: issue reading string from file: %s\n", scanner.Text())
				} else {
					l[strings.Trim(split[1], " ")] = strings.Trim(split[0], " ")
					link <- l
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error with scanner: %s\n", scanner.Text())
			}
		}
		close(link)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for js := range results {
		outputHandler(js)
	}

	//fmt.Println("out of list handler")

}

func getJSON(link <- chan map[string]string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for m := range link {
		for k, v := range m {
			results <- getJSONFromLocal(k, v)
		}
	}
}

// retrieve a JSON output from HTTPreserve without talking to the server
func httpreserveJSONOutput(link string, filename string) string {
	return getJSONFromLocal(link, filename)
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
