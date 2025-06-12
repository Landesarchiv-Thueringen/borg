package main

import (
	"fmt"
	"lath/borg/internal"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	VERSION          = "2.0.0"
	DEFAULT_RESPONSE = "Borg server version %s is running"
	FILE_STORE_PATH  = "/borg/file-store"
)

type fileAnalysis struct {
	// Summary describes the overall verification result.
	Summary internal.Summary `json:"summary"`
	// Merged feature sets ...
	FeatureSets []internal.FeatureSet `json:"featureSets"`
	// ToolResults is a list of complete responses from all tools, mapped by
	// tool name.
	ToolResults []internal.ToolResult `json:"toolResults"`
	// DurationInMs represents the duration of the analysis in milliseconds.
	DurationInMs int64 `json:"durationInMs"`
}

func main() {
	log.Printf(DEFAULT_RESPONSE, VERSION)
	initServer()
	router := gin.Default()
	router.MaxMultipartMemory = 5000 << 20 // 5 GiB
	// Allow cors to integrate Borg in other applications.
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"*"})
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	corsConfig.AllowMethods = []string{"GET", "POST"}
	// It's important that the cors configuration is used before declaring the routes.
	router.Use(cors.New(corsConfig))
	router.GET("api", getDefaultResponse)
	router.GET("api/version", getVersion)
	router.POST("api/analyze", analyzeFile)
	router.Run()
}

func initServer() {
	internal.ParseConfig()
}

func getDefaultResponse(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf(DEFAULT_RESPONSE, VERSION))
}

func getVersion(c *gin.Context) {
	c.String(http.StatusOK, VERSION)
}

func analyzeFile(c *gin.Context) {
	start := time.Now()
	file, err := c.FormFile("file")
	// no file received
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "no file received",
		})
		return
	}
	// generate unique file name for storing
	filename := uuid.New().String() + "_" + file.Filename
	fileStorePath := filepath.Join(FILE_STORE_PATH, filename)
	err = c.SaveUploadedFile(file, fileStorePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "unable to save file",
		})
		return
	}
	defer os.Remove(fileStorePath)
	identResults := internal.RunIdentificationTools(filename)
	triggeredResults := internal.RunTriggeredTools(filename, identResults)
	toolResults := internal.CombineToolResults(identResults, triggeredResults)
	mergedSets := internal.MergeFeatureSets(toolResults)
	if len(mergedSets) == 0 {
		mergedSets = make([]internal.FeatureSet, 0)
	}
	tr := internal.GetSortedToolResults(identResults, triggeredResults)
	fileAnalysis := fileAnalysis{
		Summary:      internal.GetSummary(mergedSets, tr),
		FeatureSets:  mergedSets,
		ToolResults:  tr,
		DurationInMs: time.Since(start).Milliseconds(),
	}
	c.JSON(http.StatusOK, fileAnalysis)
}
