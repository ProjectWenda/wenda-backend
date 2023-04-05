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
	"github.com/google/uuid"
)

var API_ENDPOINT string = "https://discord.com/api/v10/oauth2/token"

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
		"redirect_uri":  {os.Getenv("FRONTEND_URL")},
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
		log.Printf("Failed to add user to db %s", err)
	}

	// Redirect back to frontend
	c.IndentedJSON(http.StatusOK, gin.H{"authuid": authuid})
}

func BotAuth(c *gin.Context) {
	bot_uid, discordID, discord_name := c.Query("botUID"), c.Query("discordID"), c.Query("discordName")
	if bot_uid != os.Getenv("BOT_UID") {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "invalid bot uid"})
		return
	}

	uid, err := db.GetUID(discordID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to get uid of given user"})
		return
	}

	// User is not in table
	if uid == "" {
		uid = utils.HashToken(uuid.New().String())
		db.AddUser(db.User{
			UID:         uid,
			Token:       "",
			DiscordID:   discordID,
			DiscordName: discord_name,
		})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"authuid": uid})
}
