// Package cookie provides cookie setting functionality.
package cookie

import (
	"net/http"
	"time"
)

// SetDownloadInitiatedCookie creates an http.Cookie object to notify the client of download initiation.
func SetDownloadInitiatedCookie() http.Cookie {
	expiration := time.Now().Add(1 * time.Second)
	cookie := http.Cookie{Name: "downloadStarted", Value: "1", Expires: expiration, Path: "/"}
	return cookie
}

// SetConformityFailedCookie creates an http.Cookie object to notify the client of uploaded file's bad format.
func SetConformityFailedCookie() http.Cookie {
	expiration := time.Now().Add(1 * time.Second)
	cookie := http.Cookie{Name: "conformityFailed", Value: "1", Expires: expiration, Path: "/"}
	return cookie
}
