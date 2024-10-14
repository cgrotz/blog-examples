package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var outPath = flag.String("o", "", "path to archive file the ecourager should encourage")
	flag.Parse()
	srcPaths := flag.Args()
	for _, srcPath := range srcPaths {
		//_, filename := path.Split(srcPath)
		//nameWithoutExtension := strings.TrimSuffix(filename, path.Ext(filename))
		//encourage(srcPath, fmt.Sprintf("%s%s", path.Join(*outPath, nameWithoutExtension), ".a"))
		encourage(srcPath, *outPath)
	}
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
	_, err = outputFile.WriteString("Hello World")
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	fmt.Println("File copied and appended successfully!")
}
