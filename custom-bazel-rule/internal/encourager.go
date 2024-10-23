// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	for _, arg := range os.Args {
		if strings.Contains(arg, ":") {
			srcPath, outPath := split_arg(arg)
			encourage(srcPath, outPath)
		}
	}
}

func split_arg(arg string) (string, string) {
	slice := strings.Split(arg, ":")
	return slice[0], slice[1]
}

func encourage(sourceFile string, destinationFile string) {
	inputFile, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	// Create the destination file for writing
	outputFile, err := os.Create(destinationFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		fmt.Println("Error copying file contents:", err)
		return
	}

	// Append "Hello World" to the destination file
	_, err = outputFile.WriteString("// Wohoo! Amazing Code.")
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}
}
