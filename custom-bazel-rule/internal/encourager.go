package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// We exepect that the arguments given to the encourager are tuples of
// input_path:output_path.
func main() {
	for _, arg := range os.Args {
		if strings.Contains(arg, ":") {
			srcPath, outPath := splitArg(arg)
			encourage(srcPath, outPath)
		}
	}
}

// splitArg splits a given string with a : into the part before and after the
// colon.
func splitArg(arg string) (string, string) {
	slice := strings.Split(arg, ":")
	return slice[0], slice[1]
}

// Load the sourceFile, add a encouraging message to the end, and
// save the updated file under destinationFile.
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
	_, err = outputFile.WriteString("Hello World")
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	fmt.Println("File copied and appended successfully!")
}
