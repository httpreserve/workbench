package main

import (
	"github.com/exponential-decay/httpreserve"
)

// This structure is used to communicate with the server
// we may also use some static storage in the form of Bolt DB
// the final signal to the webapp will be a empty payload
// but with the complete flag set to true so that we know there
// is no more work to be processed. ls contains a link stat
// data structure if we can recreate one from the JSON we receive
// else the js variable will contain a single JSON document to
// be processed.
type processLog struct {
	complete bool
	ls       httpreserve.LinkStats
	js       string
}
