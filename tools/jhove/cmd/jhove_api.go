package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolOutput        string
	ExtractedFeatures map[string]string
}

type TikaOutput struct {
	MimeType *string `json:"Content-Type"`
	Encoding *string `json:"Content-Encoding"`
	Size     *string `json:"Content-Length"`
}

var defaultResponse = "JHOVE API is running"
var workDir = "/borg/tools/jhove"
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
	router.GET("validate/:module", validateFile)
	addr := "0.0.0.0:" + os.Getenv("JHOVE_API_CONTAINER_PORT")
	router.Run(addr)
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func validateFile(context *gin.Context) {
	module := context.Param("module")
	if module == "" {
		errorMessage := "no JHOVE module declared"
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
		"./jhove/jhove",
		"-m",
		module+"-hul",
		"-h",
		"json",
		fileStorePath,
	)
	log.Println(cmd.String())
	jhoveOutput, err := cmd.Output()
	if err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error executing JHOVE command",
		})
		return
	}
	jhoveOutputString := string(jhoveOutput)
	log.Println(jhoveOutputString)
}

// func extractMetadata(context *gin.Context) {
// 	fileStorePath := filepath.Join(storeDir, context.Query("path"))
// 	_, err := os.Stat(fileStorePath)
// 	if err != nil {
// 		log.Println(err)
// 		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
// 			"message": "error processing file: " + fileStorePath,
// 		})
// 		return
// 	}
// 	log.Println(fileStorePath)
// 	cmd := exec.Command(
// 		"java",
// 		"-jar",
// 		filepath.Join(workDir, "bin/tika-app-2.9.0.jar"),
// 		"--metadata",
// 		"--json",
// 		fileStorePath,
// 	)
// 	log.Println(cmd.String())
// 	tikaOutput, err := cmd.Output()
// 	if err != nil {
// 		log.Println(err)
// 		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"message": "error executing Tika command",
// 		})
// 		return
// 	}
// 	tikaOutputString := string(tikaOutput)
// 	processTikaOutput(context, tikaOutputString)
// }

// func processTikaOutput(context *gin.Context, output string) {
// 	var parsedTikaOutput TikaOutput
// 	err := json.NewDecoder(strings.NewReader(output)).Decode(&parsedTikaOutput)
// 	if err != nil {
// 		log.Println(err)
// 		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"message": "unable parse Tika output",
// 		})
// 		return
// 	}
// 	extractedFeatures := make(map[string]string)
// 	response := ToolResponse{
// 		ToolOutput: output,
// 	}
// 	if parsedTikaOutput.MimeType != nil {
// 		extractedFeatures["mimeType"] = *parsedTikaOutput.MimeType
// 	}
// 	if parsedTikaOutput.Encoding != nil {
// 		extractedFeatures["encoding"] = *parsedTikaOutput.Encoding
// 	}
// 	if parsedTikaOutput.Size != nil {
// 		extractedFeatures["size"] = *parsedTikaOutput.Size
// 	}
// 	response.ExtractedFeatures = extractedFeatures
// 	context.JSON(http.StatusOK, response)
// }
