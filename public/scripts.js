/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

function focus() {
    document.getElementById("shell-input").focus();
}

function sendCmd() {
	var xhttp = new XMLHttpRequest();
	outputElem = document.getElementById("shell-output");
	inputElem = document.getElementById("shell-input");
	xhttp.onreadystatechange = function() {
	    if (this.readyState == 4 && this.status == 200) {
			curText = outputElem.innerHTML;
	        outputElem.innerHTML = curText + this.responseText + "<br/>";
			outputElem.scrollTop = outputElem.scrollHeight;
			inputElem.value = "";
		}
	};
	xhttp.open("POST", "/ajax", true);
	text = inputElem.value;
	xhttp.send(text);
	return false;
}
