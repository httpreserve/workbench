package main

import (
	"encoding/json"
	"fmt"
	"github.com/exponential-decay/httpreserve"
	kval "github.com/kval-access-language/kval-boltdb"
	"github.com/speps/go-hashids"
	"os"
	"time"
)

// values to use to create hashid
var salt = "httpreserve"
var namelen = 8

// bucket constants
const linkIndex = "link index"
const fnameIndex = "filename index"

func getName() []int {
	t := time.Now()
	i1 := t.Minute()
	i2 := t.Second()
	return []int{i1, i2}
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

func convertInterface(v interface{}) string {
	var val string
	switch v.(type) {
	case string:
		val = fmt.Sprintf("%s", v)
	case int:
		val = fmt.Sprintf("%d", v)
	}

	if val == "" {
		return "\"\""
	}
	return val
}

func makeLinkIndex(kb kval.Kvalboltdb, lmap map[string]interface{}) {
	for k, v := range lmap {
		_, err := kval.Query(kb, "INS "+linkIndex+" >>>> "+k+" :: "+convertInterface(v))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func makeFnameIndex(kb kval.Kvalboltdb, lmap map[string]interface{}) {
	for k, v := range lmap {
		_, err := kval.Query(kb, "INS "+fnameIndex+" >> "+convertInterface(lmap["filename"])+" >>>> "+k+" :: "+convertInterface(v))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

const boltdir = "db/"

func makeBoltDir() {
	if _, err := os.Stat(boltdir); os.IsNotExist(err) {
		err := os.Mkdir(boltdir, 0700)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func boltdbHandler(ch chan string) {
	boltname := configureHashID()
	makeBoltDir()

	kb, err := kval.Connect(boltdir + "HP_" + boltname + ".bolt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening bolt database: %+v\n", err)
		os.Exit(1)
	}
	defer kval.Disconnect(kb)

	var ls httpreserve.LinkStats

	for range linkmap {
		ce := <-ch
		fmt.Println(ce)
		err := json.Unmarshal([]byte(ce), &ls)
		if err != nil {
			fmt.Fprintln(os.Stderr, "problem unmarshalling data.", err)
		}
		lmap := storeStruct(ls)
		makeLinkIndex(kb, lmap)
		makeFnameIndex(kb, lmap)
	}
}
