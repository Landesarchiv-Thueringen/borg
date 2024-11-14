package main

import (
	"lath/borg/internal"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fileAnalysis struct {
	// Summary describes the overall verification result.
	Summary internal.Summary `json:"summary"`
	// Features is an accumulated and weighted list of feature values that has
	// been extracted from different tools.
	//
	// Each feature is mapped by a key (e.g. "mimeType") and has an array of
	// values (e.g. "application/pdf") associated with a score value between 0
	// and 1 and a list of tools supporting this value. The list of values is
	// sorted by score in descending order, i.e., highest score first.
	Features map[string][]internal.FeatureValue `json:"features"`
	// ToolResults is a list of complete responses from all tools, mapped by
	// tool name.
	ToolResults []internal.ToolResult `json:"toolResults"`
}

const version = "1.3.0"
const defaultResponse = "borg server is running"
const storePath = "/borg/file-store"

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 1000 << 20 // 1 GiB
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"*"})
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	corsConfig.AllowMethods = []string{"GET", "POST"}
	// It's important that the cors configuration is used before declaring the routes.
	router.Use(cors.New(corsConfig))
	router.GET("", getDefaultResponse)
	router.GET("api/version", getVersion)
	router.POST("api/analyze-file", analyzeFile)
	router.Run()
}

func getDefaultResponse(c *gin.Context) {
	c.String(http.StatusOK, defaultResponse)
}

func getVersion(c *gin.Context) {
	c.String(http.StatusOK, version)
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
	fileStorePath := filepath.Join(storePath, filename)
	err = c.SaveUploadedFile(file, fileStorePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "unable to save file",
		})
		return
	}
	defer os.Remove(fileStorePath)
	results := internal.RunIdentificationTools(filename)
	results = internal.RunValidationTools(filename, results)
	features := internal.AccumulateFeatures(results)
	fileAnalysis := fileAnalysis{
		Summary:     internal.GetSummary(features, results),
		Features:    features,
		ToolResults: results,
	}
	c.JSON(http.StatusOK, fileAnalysis)
}
