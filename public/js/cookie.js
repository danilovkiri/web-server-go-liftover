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
        document.getElementById('dimmer').style.display = 'none';
        document.getElementById("canvas").style = 'transform: rotate(0deg);';
        document.getElementById("wrapper").style.animation = 'None';
        document.getElementById("exhaust-flame").style = 'visibility: hidden;';
        document.getElementById("exhaust-fumes-1").style = 'visibility: hidden;';
        document.getElementById("exhaust-fumes-2").style = 'visibility: hidden;';
        document.getElementById("exhaust-fumes-3").style = 'visibility: hidden;';
        document.getElementById("exhaust-fumes-4").style = 'visibility: hidden;';
        document.getElementById("exhaust-fumes-5").style = 'visibility: hidden;';
        document.getElementById("line1").style = 'visibility: hidden;';
        document.getElementById("line2").style = 'visibility: hidden;';
        document.getElementById("line3").style = 'visibility: hidden;';
        document.getElementById("line4").style = 'visibility: hidden;';
        removeCustomAlert();
        console.log('Cookie downloadStarted=1 was retrieved and screen must become active');
    } else if (getCookie("conformityFailed") === "1") {
        document.getElementById('dimmer').style.display = 'none';
        document.getElementById("canvas").style = 'transform: rotate(0deg);';
        document.getElementById("wrapper").style.animation = 'none';
        document.getElementById("exhaust-flame").style = 'visibility: hidden;';
        document.getElementById("exhaust-fumes-1").style = 'visibility: hidden;';
        document.getElementById("exhaust-fumes-2").style = 'visibility: hidden;';
        document.getElementById("exhaust-fumes-3").style = 'visibility: hidden;';
        document.getElementById("exhaust-fumes-4").style = 'visibility: hidden;';
        document.getElementById("exhaust-fumes-5").style = 'visibility: hidden;';
        document.getElementById("line1").style = 'visibility: hidden;';
        document.getElementById("line2").style = 'visibility: hidden;';
        document.getElementById("line3").style = 'visibility: hidden;';
        document.getElementById("line4").style = 'visibility: hidden;';
        removeCustomAlert();
        console.log("22", document.getElementById("canvas").style.transform)
        console.log('Cookie conformityFailed=1 was retrieved and screen must become active');
        createCustomAlert('Oops! The provided file did not meet the conformity criteria. Please, check your file and try again.');

    } else {
        downloadTimeout = setTimeout(checkDownloadCookie, 500); //Re-run this function
    }
};

const cookieSettle = function() {
    console.log('cookieSettle initiated');
    setCookie('downloadStarted', 0, 60 * 1000);
    checkDownloadCookie();
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function fallAsleep(period) {
    await sleep(period);
}

const relaxWindow = function() {
    console.log('Got a change in HTTP response iframe, relaxing screen now');
    document.getElementById('dimmer').style.display='none';
    fallAsleep(500).then(r => {console.log('Relaxing timeout now');clearTimeout(downloadTimeout)})
}

