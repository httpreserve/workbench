package main

// list handler to help us kick off some go channels
// we pass a first class function to help route our output
func listHandler(outputHandler func(ch chan string)) {
	ch := make(chan string)
	for l, f := range linkmap {
		go channelLocalLink(l, f, ch)
	}
	outputHandler(ch)
}
