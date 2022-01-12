package cookie

import (
	"net/http"
	"time"
)

func SetDownloadInitiatedCookie() http.Cookie {
	expiration := time.Now().Add(5 * time.Second)
	cookie := http.Cookie{Name: "downloadStarted", Value: "1", Expires: expiration, Path: "/"}
	return cookie
}
