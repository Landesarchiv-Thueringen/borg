/* BorgFormat - File format identification and validation
 * Copyright (C) 2024 Landesarchiv Thüringen
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type ToolResponse struct {
	ToolOutput        *string
	OutputFormat      *string
	ExtractedFeatures *map[string]string
	Error             *string
}

type TikaOutput struct {
	MimeType    *string `json:"Content-Type"`
	Encoding    *string `json:"Content-Encoding"`
	PDFVersion  *string `json:"pdf:PDFVersion"`
	PDFAVersion *string `json:"pdfa:PDFVersion"`
}

var defaultResponse = "Tika API is running"
var workDir = "/borg/tools/tika"
var storeDir = "/borg/file-store"
var outputFormat = "json"

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/extract-metadata", extractMetadata)
	router.Run("0.0.0.0:80")
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func extractMetadata(context *gin.Context) {
	fileStorePath := filepath.Join(storeDir, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		errorMessage := "error processing file: " + fileStorePath
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	cmd := exec.Command(
		"java",
		"-jar",
		filepath.Join(workDir, "third_party/tika-app-2.9.0.jar"),
		"--metadata",
		"--"+outputFormat,
		fileStorePath,
	)
	tikaOutput, err := cmd.Output()
	if err != nil {
		errorMessage := "error executing Tika command"
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	tikaOutputString := string(tikaOutput)
	processTikaOutput(context, tikaOutputString)
}

func processTikaOutput(context *gin.Context, output string) {
	var parsedTikaOutput TikaOutput
	err := json.NewDecoder(strings.NewReader(output)).Decode(&parsedTikaOutput)
	if err != nil {
		errorMessage := "unable parse Tika output"
		log.Println(errorMessage)
		log.Println(err)
		response := ToolResponse{
			ToolOutput:   &output,
			OutputFormat: &outputFormat,
			Error:        &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	extractedFeatures := make(map[string]string)
	response := ToolResponse{
		ToolOutput:        &output,
		OutputFormat:      &outputFormat,
		ExtractedFeatures: &extractedFeatures,
	}
	if parsedTikaOutput.MimeType != nil {
		// removes charset from MIME-Type if existing, example: text/x-yaml; charset=ISO-8859-1
		mimeType := strings.Split(*parsedTikaOutput.MimeType, ";")[0]
		extractedFeatures["mimeType"] = mimeType
	}
	if parsedTikaOutput.Encoding != nil {
		extractedFeatures["encoding"] = *parsedTikaOutput.Encoding
	}
	// use PDF/A version if existing
	if parsedTikaOutput.PDFAVersion != nil {
		extractedFeatures["formatVersion"] = "PDF/" + *parsedTikaOutput.PDFAVersion
	} else if parsedTikaOutput.PDFVersion != nil {
		// no PDF/A version --> use normal version info
		extractedFeatures["formatVersion"] = *parsedTikaOutput.PDFVersion
	}
	context.JSON(http.StatusOK, response)
}
