package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	api := r.Group("/api")
	api.POST("/test", func(c *gin.Context) {
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
		fmt.Println(res)
		c.AbortWithStatus(http.StatusOK)
	})
	r.Run(":8082")
}
