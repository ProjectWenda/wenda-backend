package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

var API_ENDPOINT string = "https://discord.com/api/v10/oauth2/token"
var REDIRECT_URI string = "http://localhost:3000"

func GetAuth(c *gin.Context) {
	// Get params for auth
	client_id, client_secret := os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")
	state := c.Query("state")
	code := c.Query("code")

	fmt.Println("Code is", code)
	fmt.Println("State is", state)
	fmt.Println("Client id is", client_id)
	fmt.Println("Client secret is", client_secret)

	// Form auth URL
	data := url.Values{
		"client_id":     {client_id},
		"client_secret": {client_secret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {REDIRECT_URI},
	}

	fmt.Println(data.Encode())

	// resp, err := http.PostForm(API_ENDPOINT, data)

	// if err != nil {
	// 	fmt.Println("Issue posting to discord auth api")
	// }

	// var res map[string]interface{}

	// json.NewDecoder(resp.Body).Decode(&res)

	c.String(http.StatusOK, "hi %s", "bob")
}
