// function to save the page to the internet archive
// but let httpreserve manage the transaction such that 
// the user doesn't have to leave the page. 
// Simple AJAX function to make our demo page nice and clean
function saveToInternetArchive(saveLink) {

	var iaSaved = "httpreserve-saved";
	var httpreserveError = "httpreserve-error"; 
   var method = "post";

	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState == XMLHttpRequest.DONE ) {
			if (xmlhttp.status == 200) {
				document.getElementById(iaSaved).innerHTML = xmlhttp.responseText;
			}
			else if (xmlhttp.status == 400) {
            document.getElementById(httpreserveError).innerHTML = '<div class="error"><br/>[IA SAVE WARNING] There was an error 400<br/></div><br/>';
			}
			else {
            document.getElementById(httpreserveError).innerHTML = '<div class="error"><br/>[IA SAVE WARNING] something else other than 200 was returned<br/><br/></div><br/>';
			}
		}
	};

	var key="saveurl";
	var value = saveLink;

   if (method.toLowerCase() == "post") {
      xmlhttp.open("POST", "save", true);
      xmlhttp.setRequestHeader("Content-type","application/x-www-form-urlencoded");
      xmlhttp.send(key + "=" + value);	// Placeholder, not yet implemented
      return;
   }
   
   if (method.toLowerCase() == "get") {
      xmlhttp.open("GET", "save?" + key + "=" + value, true);
      xmlhttp.send(key + "=" + value);	// Placeholder, not yet implemented   
      return;
   }
   
   document.getElementById("httpreserve-analysis").innerHTML = '[WARNING] issue parsing the form in JavaScript';
}
