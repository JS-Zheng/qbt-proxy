package main

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Creates a Gin handler function that proxies incoming HTTP requests to a qBittorrent instance.
func ApiProxy(config Config) gin.HandlerFunc {
	proxy := createReverseProxy(config)

	return func(c *gin.Context) {
		err := handleSID(c)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Failed to handle SID", err)
			return
		}

		// Forward the request to the target qBittorrent Web API server
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// Creates a reverse proxy using the provided configuration
func createReverseProxy(config Config) *httputil.ReverseProxy {
	targetURL, err := url.Parse(config.BaseURL)
	if err != nil {
		panic("Invalid base URL")
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	oriDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		oriDirector(req)
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		// Set the request's Host header to match the target server
		req.Host = targetURL.Host
        // qBittorrent cannot handle them well
		req.Header.Del("Origin")
		req.Header.Del("Referer")
	}
	return proxy
}

// Checks if the SID is already in the request header, extracts it if necessary,
// sets the SID as a cookie header, and removes the SID from the URL if it exists.
func handleSID(c *gin.Context) error {
	// Check if the SID is already in the request header
	if _, err := c.Request.Cookie("SID"); err != nil {
		// Extract SID from request URL query parameter if not in the header
		sid, err := extractSIDFromURL(c.Request)
		if err != nil {
			return err
		}
		setSIDToCookie(c.Request, sid)
	}

	if sidExistsInURL(c.Request) {
		removeSIDFromURL(c.Request)
	}

	return nil
}

// Extracts the SID from the URL's query parameters
func extractSIDFromURL(r *http.Request) (string, error) {
	sid := r.URL.Query().Get("sid")
	if sid == "" {
		return "", errors.New("SID not found in request URL query parameter")
	}
	return sid, nil
}

// Checks if the SID exists in the URL's query parameters
func sidExistsInURL(r *http.Request) bool {
	return r.URL.Query().Get("sid") != ""
}

// Removes the SID query parameter from the URL
func removeSIDFromURL(r *http.Request) {
	q := r.URL.Query()
	q.Del("sid")
	r.URL.RawQuery = q.Encode()
}

// Sets the SID as a cookie header in the request
func setSIDToCookie(r *http.Request, sid string) {
	// Set SID as cookie header
	cookie := &http.Cookie{Name: "SID", Value: sid}
	r.AddCookie(cookie)
}
