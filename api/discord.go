package api

import (
	"app/wenda/db"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var DISCORD_BASE string = "https://discord.com/api/v10"

type AuthResponse struct {
	Application interface{} `json:"application"`
	Scopes      []string    `json:"scopes"`
	Expires     string      `json:"expires"`
	User        struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Avatar        string `json:"avatar"`
		Discriminator string `json:"discriminator"`
		PublicFlags   int    `json:"public_flags"`
	} `json:"user"`
}

func DiscordAuthData(token string) (string, string) {
	req, err := http.NewRequest("GET", DISCORD_BASE+"/oauth2/@me", nil)
	if err != nil {
		fmt.Println("Failed to form discord authentication info req")
		return "", ""
	}
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to get discord authentication info")
		return "", ""
	}

	var res AuthResponse
	json.NewDecoder(resp.Body).Decode(&res)

	return res.User.ID, res.User.Username
}

func GetDiscordUser(c *gin.Context) {
	uid := c.Query("uid")
	token, err := db.GetUserToken(uid)
	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "failed to retrieve user token"})
		return
	}

	req, err := http.NewRequest("GET", DISCORD_BASE+"/users/@me", nil)
	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "failed to form user get req"})
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "failed to fetch user data"})
		return
	}

	// Parse token into res
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	c.IndentedJSON(http.StatusOK, res)
}

func GetDiscordFriends(c *gin.Context) {
	uid := c.Query("uid")
	token := User_id_token[uid]

	req, err := http.NewRequest("GET", DISCORD_BASE+"/users/@me/relationships", nil)
	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "failed to form user friends get req"})
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "failed to fetch user friends"})
		return
	}

	// Parse token into res
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	c.IndentedJSON(http.StatusOK, res)
}
