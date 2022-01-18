const setCookie = function(name, value, expiracy) {
    var exdate = new Date();
    exdate.setTime(exdate.getTime() + expiracy);
    var c_value = escape(value) + ((expiracy == null) ? "" : "; expires=" + exdate.toUTCString());
    document.cookie = name + "=" + c_value + '; path=/';
    console.log(`Cookie ${name} was set with value ${value}`)
};

const getCookie = function(name) {
    var i, x, y, ARRcookies = document.cookie.split(";");
    console.log("Found cookies: ", document.cookie.split(";"))
    for (i = 0; i < ARRcookies.length; i++) {
        x = ARRcookies[i].substr(0, ARRcookies[i].indexOf("="));
        y = ARRcookies[i].substr(ARRcookies[i].indexOf("=") + 1);
        x = x.replace(/^\s+|\s+$/g, "");
        if (x == name) {
            return y ? decodeURI(unescape(y.replace(/\+/g, ' '))) : y; //;//unescape(decodeURI(y));
        }
    }
};

let downloadTimeout;
const checkDownloadCookie = function() {
    if (getCookie("downloadStarted") === "1") {
        setCookie("downloadStarted", 1, 1000); //Expiration could be anything... As long as we reset the value
        document.getElementById('dimmer').style.display='none';
        removeCustomAlert();
        console.log('Cookie downloadStarted=1 was set and screen must become active')
    } else if (getCookie("conformityFailed") === "1") {
        setCookie("downloadStarted", 0, 1000);
        document.getElementById('dimmer').style.display='none';
        removeCustomAlert();
        console.log('Cookie downloadStarted=0 was set and screen must become active')
        createCustomAlert('Oops! The provided file did not meet the conformity criteria. Please, check your file and try again.')
    } else {
        downloadTimeout = setTimeout(checkDownloadCookie, 500); //Re-run this function in 1 second.
    }
};

const cookieSettle = function() {
    console.log('cookieSettle initiated')
    setCookie('downloadStarted', 0, 60 * 1000);
    setTimeout(checkDownloadCookie, 500);
}

const relaxWindow = function() {
    console.log('Got an error in iframe, relaxing screen now')
    document.getElementById('dimmer').style.display='none';
    clearTimeout(downloadTimeout)
}

