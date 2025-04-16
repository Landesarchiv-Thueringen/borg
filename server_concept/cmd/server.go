package main

import (
	"fmt"
	"lath/borg/internal"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	VERSION          = "2.0.0"
	DEFAULT_RESPONSE = "Borg server version %s is running"
	FILE_STORE_PATH  = "/borg/file-store"
)

var serverConfig internal.ServerConfig

func main() {
	log.Printf(DEFAULT_RESPONSE, VERSION)
	initServer()
	router := gin.Default()
	router.MaxMultipartMemory = 5000 << 20 // 5 GiB
	router.SetTrustedProxies([]string{})
	router.GET("api", getDefaultResponse)
	router.GET("api/version", getVersion)
	router.Run()
}

func initServer() {
	serverConfig = internal.ParseConfig()
}

func getDefaultResponse(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf(DEFAULT_RESPONSE, VERSION))
}

func getVersion(c *gin.Context) {
	c.String(http.StatusOK, VERSION)
}
