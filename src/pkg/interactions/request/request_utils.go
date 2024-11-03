package request

import (
	"net/http"
	"strings"
)

func IsMultiPart(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data")
}

func IsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

func IsXML(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/xml")
}

func IsForm(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
}

func IsText(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "text/plain")
}

func IsHTML(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "text/html")
}
