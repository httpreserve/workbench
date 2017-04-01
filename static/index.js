// Auto request information from the server...
var time = new Date().getTime();

function refresh() {
	if (timer > 0) {
		if(new Date().getTime() - time >= timer) {
			httpreserve("post");
			setTimeout(refresh, timer);
		}
	}
}

var timer = 1000;
setTimeout(refresh, timer);

// Simple AJAX function to make our demo page nice and clean
function httpreserve(formMethod) {

	var analysistable = "httpreserve-analysis";
   var method = formMethod;

	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState == XMLHttpRequest.DONE ) {
			if (xmlhttp.status == 200) {
				document.getElementById(analysistable).innerHTML = formatRow(xmlhttp.responseText);
			}
			else if (xmlhttp.status == 400) {
            document.getElementById(analysistable).innerHTML = '[WARNING] There was an error 400';
			}
			else {
            document.getElementById(analysistable).innerHTML = '[WARNING] something else other than 200 was returned';
			}
		}
	};

   if (method.toLowerCase() == "post") {
      xmlhttp.open("POST", "httpreserve", true);
      xmlhttp.setRequestHeader("Content-type","application/x-www-form-urlencoded");
      xmlhttp.send("key" + "=" + "value");	// Placeholder, not yet implemented
      return;
   }
   
   if (method.toLowerCase() == "get") {
      xmlhttp.open("GET", "httpreserve?" + key + "=" + value, true);
      xmlhttp.send("key" + "=" + "value");	// Placeholder, not yet implemented   
      return;
   }
   
   document.getElementById("httpreserve-analysis").innerHTML = '[WARNING] issue parsing the form in JavaScript';
}

// Global for all our table data
// TODO: monitor performance as it may contain a lot of data
var content = ""

// Function to fomrat the data output by the server when it is received
function formatRow(data_arr) {

	var tableStart = "<table><th>time</th>"
	var trStart = "<tr><td>"
	var trEnd = "</td></tr>"
	var tableEnd = "</table>"
	var padding = "<br/><br/><br/><br/>"

	var arr = data_arr.split(",", 2)

	if (data_arr.length > 0) {	
		newRow = trStart + arr[1] + trEnd;
	} else {
		newRow = trStart + "server data issue" + trEnd;
	}

	//if (arr[0] == "false") {
	//	timer = 0
	//}

	content = content + newRow
	return tableStart + content + tableEnd + padding;	
}

