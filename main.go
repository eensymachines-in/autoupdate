package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

var (
	REPO_NAME       = "" // expected repository name
	REPO_DIR_ONHOST = "" // directory on the host where the git repo is located
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

func init() {
	log.SetFormatter(&log.TextFormatter{DisableColors: false, FullTimestamp: false})
	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)

	REPO_NAME = os.Getenv("REPO_NAME")
	if REPO_NAME == "" {
		log.Panic("Empty environment variable: REPO_NAME")
	}
	REPO_DIR_ONHOST = os.Getenv("REPO_DIR")
	if REPO_NAME == "" {
		log.Panic("Empty environment variable: REPO_DIR_ONHOST")
	}
}

func main() {
	log.Info("=========")
	log.Info("Starting the auto deploy service..")
	log.Info("=========")

	defer log.Warn("Now closing the patio-web program...")

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
		if REPO_NAME != res.Repository.Name {
			log.WithFields(log.Fields{
				"expected": REPO_NAME,
				"got":      res.Repository.Name,
			}).Error("Repository name isnt as expected, did you change the repository name upstream?")
			c.AbortWithStatus(http.StatusOK) //but send back 200 ok to the server, acknowledge
			return
		}
		if res.Action != "published" {
			log.WithFields(log.Fields{
				"expected": "published",
				"got":      res.Action,
			}).Error("Is only scheduled to run when new release is created")
			c.AbortWithStatus(http.StatusOK) //but send back 200 ok to the server, acknowledge
			return
		}
		// changing the current directory before executing the other bash scripts
		if err := os.Chdir(REPO_DIR_ONHOST); err != nil {
			log.WithFields(log.Fields{
				"expected": REPO_DIR_ONHOST,
				"err":      err,
			}).Error("Error changing to the working directory on the host")
			c.AbortWithStatus(http.StatusOK) //but send back 200 ok to the server, acknowledge
			return
		}
		log.Info("Changed to directory..")
		fi, err := os.ReadDir(REPO_DIR_ONHOST)
		if err != nil {
			log.WithFields(log.Fields{
				"expected": "file information for the repo directory",
				"err":      err,
			}).Error("Error changing to the working directory on the host")
			c.AbortWithStatus(http.StatusOK) //but send back 200 ok to the server, acknowledge
			return
		}
		for _, entry := range fi {
			log.Debug(entry.Name())
		}
		c.AbortWithStatus(http.StatusOK)
	})
	r.Run(":8082")
}
