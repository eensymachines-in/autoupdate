package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReleaseInfo struct {
	Body       string `json:"body"`
	CreatedAt  string `json:"created_at"`
	Draft      bool   `json:"draft"`
	PreRelease bool   `json:"prerelease"`
	Tag        string `json:"tag_name"`
}
type RepoInfo struct {
	Name          string `json:"name"`
	CloneUrl      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
}

type SenderInfo struct {
}

type WebHkRelease struct { // payload we receive from web hook notification
	Action     string      `json:"action"`
	Release    ReleaseInfo `json:"release"`
	Repository RepoInfo    `json:"repository"`
	Sender     SenderInfo  `json:"sender"`
}

func CORS(c *gin.Context) {
	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")
	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		// Everytime we receive an OPTIONS request,
		// we just return an HTTP 200 Status Code
		// Like this, Angular can now do the real
		// request using any other method than OPTIONS
		c.AbortWithStatus(http.StatusOK)
	}
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	api := r.Group("/api")
	api.Use(CORS).POST("/test", func(c *gin.Context) {
		fmt.Println("We have received hook notification..")
		res := WebHkRelease{}
		err := c.ShouldBind(&res)
		defer c.Request.Body.Close()

		if err != nil {
			fmt.Println("Error unmarshaling the payload")
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		fmt.Println("received payload from github..")
		fmt.Println(res)
		c.AbortWithStatus(http.StatusOK)
	})
	r.Run(":8082")
}
