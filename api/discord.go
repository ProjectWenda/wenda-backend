package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

var DISCORD_BASE string = "https://discord.com/api/v10"

func GetDiscordUser(c *gin.Context) {
	uid := c.Query("uid")
	token := User_id_token[uid]

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
