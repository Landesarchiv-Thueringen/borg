package main

import (
	"fmt"
	"lath/borg/internal"
	"log"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"

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
}

func main() {
	log.Printf(DEFAULT_RESPONSE, VERSION)
	initServer()
	router := gin.Default()
	router.MaxMultipartMemory = 5000 << 20 // 5 GiB
	router.SetTrustedProxies([]string{})
	router.GET("api", getDefaultResponse)
	router.GET("api/version", getVersion)
	router.POST("api/analyze-file", analyzeFile)
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
	identificationResults := internal.RunIdentificationTools(filename)
	triggeredResults := internal.RunTriggeredTools(filename, identificationResults)
	toolResults := make(map[string]internal.ToolResult)
	for k, v := range identificationResults {
		toolResults[k] = v
	}
	for k, v := range triggeredResults {
		toolResults[k] = v
	}
	mergedSets := internal.MergeFeatureSets(toolResults)
	for _, s := range mergedSets {
		log.Println(s)
	}
	tr := slices.Collect(maps.Values(toolResults))
	fileAnalysis := fileAnalysis{
		Summary:     internal.GetSummary(mergedSets, tr),
		FeatureSets: mergedSets,
		ToolResults: tr,
	}
	c.JSON(http.StatusOK, fileAnalysis)
}
