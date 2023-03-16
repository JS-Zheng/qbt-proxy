package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type LoginResponse struct {
	Sid string `json:"sid"`
}

// Creates a Gin handler function that proxies a login request to a qBittorrent instance.
func LoginProxy(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract username and password from the request parameters
		username, password := extractCredentials(c)

		// Create a login request body
		data := createRequestBody(username, password)

		// Send the login request to the qBittorrent Web API
		resp, err := sendLoginRequest(config, data)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, "Failed to connect to qBittorrent", err)
			return
		}
		defer resp.Body.Close()

		// Extract the SID from the response cookie
		sid := extractSIDFromCookies(resp)

		// Respond with the SID in JSON format or an error message
		if sid != "" {
			c.JSON(http.StatusOK, LoginResponse{Sid: sid})
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		msg := fmt.Sprintf("Failed to login to qBittorrent: %s", string(body))
		HandleError(c, http.StatusBadRequest, msg, err)
	}
}

func extractCredentials(c *gin.Context) (string, string) {
	return c.PostForm("username"), c.PostForm("password")
}

func createRequestBody(username, password string) url.Values {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	return data
}

func sendLoginRequest(config Config, data url.Values) (*http.Response, error) {
	return http.PostForm(fmt.Sprintf("%s/api/v2/auth/login", config.BaseURL), data)
}

func extractSIDFromCookies(resp *http.Response) string {
	var sid string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "SID" {
			sid = cookie.Value
			break
		}
	}
	return sid
}
