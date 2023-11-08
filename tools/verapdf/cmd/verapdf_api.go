package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolOutput        string
	OutputFormat      string
	ExtractedFeatures map[string]string
}

type VeraPDFOutput struct {
	Report Report `json:"report"`
}

type Report struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	ValidationResult ValidationResult `json:"validationResult"`
}

type ValidationResult struct {
	ProfileName string `json:"profileName"`
	Compliant   bool   `json:"compliant"`
}

var defaultResponse = "veraPDF API is running"
var workDir = "/borg/tools/verapdf"
var storeDir = "/borg/filestore"

func main() {
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"*"})
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	corsConfig.AllowMethods = []string{"GET", "POST"}
	// It's important that the cors configuration is used before declaring the routes.
	router.Use(cors.New(corsConfig))
	router.GET("", getDefaultResponse)
	router.GET("/validate/:profile", validateFile)
	addr := "0.0.0.0:" + os.Getenv("VERAPDF_API_CONTAINER_PORT")
	router.Run(addr)
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func validateFile(context *gin.Context) {
	profile := context.Param("profile")
	if profile == "" {
		errorMessage := "no veraPDF profile declared"
		log.Println(errorMessage)
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": errorMessage,
		})
		return
	}
	fileStorePath := filepath.Join(storeDir, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "error processing file: " + fileStorePath,
		})
		return
	}
	cmd := exec.Command(
		"/bin/ash",
		filepath.Join(workDir, "bin/verapdf"),
		"-f",
		profile,
		"--format",
		"json",
		"-v",
		fileStorePath,
	)
	veraPDFOutput, err := cmd.Output()
	// exit status 1: file for profile invalid but validation job succesfull XD
	if err != nil && err.Error() != "exit status 1" {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error executing veraPDF command",
		})
		return
	}
	veraPDFOutputString := string(veraPDFOutput)
	processVeraPDFOutput(context, veraPDFOutputString)
}

func processVeraPDFOutput(context *gin.Context, output string) {
	var veraPDFOutput VeraPDFOutput
	err := json.NewDecoder(strings.NewReader(output)).Decode(&veraPDFOutput)
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "unable parse veraPDF output",
		})
		return
	}
	extractedFeatures := make(map[string]string)
	response := ToolResponse{
		ToolOutput:        output,
		OutputFormat:      "json",
		ExtractedFeatures: extractedFeatures,
	}
	if veraPDFOutput.Report.Jobs != nil && len(veraPDFOutput.Report.Jobs) > 0 {
		extractedFeatures["valid"] = strconv.FormatBool(
			veraPDFOutput.Report.Jobs[0].ValidationResult.Compliant,
		)
	}
	context.JSON(http.StatusOK, response)
}
