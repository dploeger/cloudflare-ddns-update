// Package internal includes internal tools used by cloudflare-ddns-update
package internal

import (
	"github.com/cloudflare/cloudflare-go/v2/option"
	"log"
	"net/http"
)

// RequestResponseLogger is a cloudflare debug logging middleware
func RequestResponseLogger(req *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
	// Before the request
	log.Println(req)

	// Forward the request to the next handler
	res, err = next(req)

	// Handle stuff after the request
	log.Println(res, err)

	return res, err
}
