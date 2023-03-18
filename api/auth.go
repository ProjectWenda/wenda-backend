package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"app/wenda/db"
	"app/wenda/utils"

	"github.com/gin-gonic/gin"
)

var API_ENDPOINT string = "https://discord.com/api/v10/oauth2/token"
var REDIRECT_URI string = "http://localhost:8080/auth"
var FRONTEND_URI string = "http://localhost:5173"

// Initialize an empty map, eventually we want to be storing this in the DB
// Maps state -> auth token
var User_id_token map[string]string = make(map[string]string)

func GetAuth(c *gin.Context) {
	// Get params for auth from url
	client_id, client_secret := os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")
	code := c.Query("code")

	// URL encode params to pass into token auth
	data := url.Values{
		"client_id":     {client_id},
		"client_secret": {client_secret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {REDIRECT_URI},
	}

	// Post to token auth
	resp, err := http.PostForm(API_ENDPOINT, data)
	if err != nil {
		fmt.Println("Issue posting to discord auth api")
		// Redirect to an error page here
	}

	// Parse token into res
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	token := res["access_token"].(string)
	// Hash the token
	authuid := utils.HashToken(token)

	discord_id, discord_name := DiscordAuthData(token)

	new_user := db.User{
		UID:         authuid,
		Token:       res["access_token"].(string),
		DiscordID:   discord_id,
		DiscordName: discord_name,
	}

	err = db.AddUser(new_user)
	if err != nil {
		log.Fatalf("Faile dto add user to db %s", err)
	}

	// Map hash -> token
	User_id_token[authuid] = res["access_token"].(string) // type to string

	// Set cookie for redirect
	authid_cookie := http.Cookie{
		Name:   "authuid",
		Value:  authuid,
		MaxAge: 604800, //1 week
		Secure: true,
	}
	http.SetCookie(c.Writer, &authid_cookie)

	// Redirect back to frontend
	c.Redirect(http.StatusMovedPermanently, FRONTEND_URI)
}
