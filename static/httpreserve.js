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

var timer = 3000;
setTimeout(refresh, timer);

// Simple AJAX function to make our demo page nice and clean
function httpreserve(formMethod) {

	var analysistable = "httpreserve-analysis";
	var httpreserveError = "httpreserve-error";
   var method = formMethod;

	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState == XMLHttpRequest.DONE ) {
			if (xmlhttp.status == 200) {
				//document.getElementById(analysistable).innerHTML = formatRow(xmlhttp.responseText);
				formatRow(xmlhttp.responseText);
			}
			else if (xmlhttp.status == 400) {
            document.getElementById(httpreserveError).innerHTML = '<div class="error"><br/>[WARNING] There was an error 400<br/></div><br/>';
			}
			else {
            document.getElementById(httpreserveError).innerHTML = '<div class="error"><br/>[WARNING] something else other than 200 was returned<br/><br/></div><br/>';
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
var content = "";
var processed = false;
var processtime = ""; 

// Function to fomrat the data output by the server when it is received
function formatRow(data_arr) {

	var httpreserveProcessing = "httpreserve-processing";

	var arr = data_arr.split("•", 2);

	if (arr[0] != "processing")
	{
		if (!processed)
		{
			//give some stats back about the processing...
		   document.getElementById(httpreserveProcessing).innerHTML = '<div class="processing"><br/>Processed: ' + processtime + '<br/></div><br/>';
			processed = true;
		}

		var tableStart = "<table><th>time</th>";
		var tableEnd = "</table>";
		var padding = "<br/><br/><br/><br/>";

		var arr = data_arr.split("•", 2);
		var newRow = "";

		if (data_arr.length > 0) {	
			newRow = arr[1];
		} else {
			timer = 0;
			return "";
		}

		if (arr[0] == "false") {
			timer = 0;
		}

		// Add a slide
		updateSlick(newRow);

		//exit function
		return
	} else {
		//give some stats back about the processing...
		processtime = arr[1];
      document.getElementById(httpreserveProcessing).innerHTML = '<div class="processing"><br/>Processing: ' + processtime + '<br/></div><br/>';
	}
}

