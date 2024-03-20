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
	"encoding/csv"
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

var defaultResponse = "DROID API is running"
var workDir = "/borg/tools/droid"
var storeDir = "/borg/file-store"
var signatureFilePath = filepath.Join(workDir, "third_party/DROID_SignatureFile_V114.xml")
var containerSignatureFilePath = filepath.Join(workDir, "third_party/container-signature-20230822.xml")
var outputFormat = "csv"

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("", getDefaultResponse)
	router.GET("/identify-file-format", identifyFileFormat)
	router.Run("0.0.0.0:80")
}

func getDefaultResponse(context *gin.Context) {
	context.String(http.StatusOK, defaultResponse)
}

func identifyFileFormat(context *gin.Context) {
	fileStorePath := filepath.Join(storeDir, context.Query("path"))
	_, err := os.Stat(fileStorePath)
	if err != nil {
		log.Println(err)
		errorMessage := "error processing file: " + fileStorePath
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	cmd := exec.Command(
		"/bin/ash",
		"/borg/tools/droid/third_party/droid.sh",
		"-Ns",
		signatureFilePath,
		"-Nc",
		containerSignatureFilePath,
		"-Nr",
		fileStorePath,
	)
	droidOutput, err := cmd.Output()
	if err != nil {
		log.Println(err)
		errorMessage := "error executing DROID command"
		response := ToolResponse{
			Error: &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	droidOutputString := string(droidOutput)
	csvReader := csv.NewReader(strings.NewReader(droidOutputString))
	formats, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err.Error())
		errorMessage := "unable to DROID csv output"
		response := ToolResponse{
			ToolOutput:   &droidOutputString,
			OutputFormat: &outputFormat,
			Error:        &errorMessage,
		}
		context.JSON(http.StatusOK, response)
		return
	}
	extractedFeatures := make(map[string]string)
	response := ToolResponse{
		ToolOutput:        &droidOutputString,
		OutputFormat:      &outputFormat,
		ExtractedFeatures: &extractedFeatures,
	}
	// TODO: discuss returning multiple results
	if len(formats) > 1 && formats[1][1] != "" {
		extractedFeatures["puid"] = formats[1][1]
	}
	context.JSON(http.StatusOK, response)
}
