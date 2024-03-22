package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		byt, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("Error reading the request payload")
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		res := map[string]interface{}{}
		err = json.Unmarshal(byt, &res)
		if err != nil {
			fmt.Println("Error unmarshaling the payload")
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		fmt.Println("received payload from github..")
		for k, v := range res {
			fmt.Printf("%s:%v\n", k, v)
		}
		c.AbortWithStatus(http.StatusOK)
	})
	r.Run(":8082")
}
