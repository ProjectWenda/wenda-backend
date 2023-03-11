package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

var API_ENDPOINT string = "https://discord.com/api/v10/oauth2/token"
var REDIRECT_URI string = "http://localhost:8080/auth"
var FRONTEND_URI string = "http://localhost:3000"

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
	}

	// Parse token into res
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	// Redirect back to frontend
	c.Redirect(http.StatusMovedPermanently, FRONTEND_URI)
}
