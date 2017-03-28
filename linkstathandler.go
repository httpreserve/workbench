package main 

import (
	"github.com/exponential-decay/httpreserve"
)
// storeStruct allows us to get a different representation of the LinkStats structure
// e.g. as a map we have good flexibility over looping and passing around without 
// reglection to iterate through the struct for us. 
func storeStruct(ls httpreserve.LinkStats) map[string]interface{} {

	var lmap = make(map[string]interface{})

	lmap["filename"] = ls.FileName
	lmap["analysis version number"] = ls.AnalysisVersionNumber
	lmap["analysis version text"] = ls.AnalysisVersionText
	lmap["link"] = ls.Link
	lmap["response code"] = ls.ResponseCode
	lmap["response text"] = ls.ResponseText
	lmap["screen shot"] = ls.ScreenShot
	lmap["internet archive latest"] = ls.InternetArchiveLinkLatest
	lmap["internet archive earliest"] = ls.InternetArchiveLinkEarliest
	lmap["internet archive save link"] = ls.InternetArchiveSaveLink
	lmap["internet archive response code"] = ls.InternetArchiveResponseCode
	lmap["internet archive response text"] = ls.InternetArchiveResponseText
	lmap["archived"] = ls.Archived
	lmap["protocol error"] = ls.ProtocolError
	lmap["protocol error"] = ls.ProtocolErrorMessage

	return lmap
}