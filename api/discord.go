package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var DISCORD_BASE string = "https://discord.com/api/v10"

func GetDiscordUser(c *gin.Context) {
	uid := c.Query("uid")
	token := User_id_token[uid]

	req, err := http.NewRequest("GET", DISCORD_BASE+"/users/@me", nil)
	if err != nil {
		fmt.Println("Failed to form get user req")
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "failed to form user get req"})
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to form get user req")
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "failed to fetch user data"})
	}

	// Parse token into res
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	c.IndentedJSON(http.StatusOK, res)
}
