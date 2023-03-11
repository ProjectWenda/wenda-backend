package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"app/wenda/utils"

	"github.com/gin-gonic/gin"
)

var API_ENDPOINT string = "https://discord.com/api/v10/oauth2/token"
var REDIRECT_URI string = "http://localhost:8080/auth"
var FRONTEND_URI string = "http://localhost:3000"

// Initialize an empty map, eventually we want to be storing this in the DB
// Maps state -> auth token
var User_id_token map[string]string = make(map[string]string)

func GetAuth(c *gin.Context) {
	// Get params for auth from url
	client_id, client_secret := os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")
	code := c.Query("code")

	fmt.Println("Code is", code)
	fmt.Println("Client ID is", client_id)
	fmt.Println("Client secret is", client_secret)

	// URL encode params to pass into token auth
	data := url.Values{
		"client_id":     {client_id},
		"client_secret": {client_secret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {REDIRECT_URI},
	}
	fmt.Println(data.Encode())

	// Post to token auth
	resp, err := http.PostForm(API_ENDPOINT, data)
	if err != nil {
		fmt.Println("Issue posting to discord auth api")
		// Redirect to an error page here
	}

	// Parse token into res
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	// Hash the token
	authuid := utils.HashToken(res["access_token"].(string))
	// Map hash -> token
	User_id_token[authuid] = res["access_token"].(string) // type to string

	// URL Encode to redirect to frontend
	frontend_params := url.Values{"authuid": {authuid}}

	// Redirect back to frontend
	c.Redirect(http.StatusMovedPermanently, FRONTEND_URI+"?"+frontend_params.Encode())
}
