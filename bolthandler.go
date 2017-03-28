package main 

import (
	"fmt"
	"github.com/speps/go-hashids"
	kval "github.com/kval-access-language/kval-boltdb"
	"os"
	"time"
)

// values to use to create hashid
var salt = "httpreserve"
var namelen = 8
//var name = []int{23, 13}

func getName() []int {
	t := time.Now()
	i1 := t.Minute()
	i2 := t.Second()
	return []int{i1,i2}
}

func configureHashID() string {
	
	name := getName()

	//hashdata
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = namelen

	//hash
	h := hashids.NewWithData(hd)
	e, _ := h.Encode(name)
	return e
}

func boltdbHandler(ch chan string) {
	boltname := configureHashID()

	kb, err := kval.Connect("HP_" + boltname + ".bolt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening bolt database: %#v", err)
		os.Exit(1)
	}
	defer kval.Disconnect(kb)	
}

/* Example JSON:
{
   "FileName": "bbc news",
   "AnalysisVersionNumber": "0.0.0",
   "AnalysisVersionText": "exponentialDK-httpreserve/0.0.0",
   "Link": "http://www.bbc.co.uk/news",
   "ResponseCode": 200,
   "ResponseText": "OK",
   "ScreenShot": "",
   "InternetArchiveLinkLatest": "http://web.archive.org/web/20170328040059/http://www.bbc.co.uk/news/",
   "InternetArchiveLinkEarliest": "http://web.archive.org/web/19971009011901/http://www.bbc.co.uk/news/",
   "InternetArchiveSaveLink": "http://web.archive.org/save/http://www.bbc.co.uk/news",
   "InternetArchiveResponseCode": 200,
   "InternetArchiveResponseText": "OK",
   "Archived": true,
   "ProtocolError": false,
   "ProtocolErrorMessage": ""
},
*/
