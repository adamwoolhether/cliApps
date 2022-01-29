package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>Markdown Preview Tool</title>
  </head>
  <body>
`
	footer = `
  </body>
</html>
`
)

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	// If no file input provided, show usage.
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(*filename, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run coordinates te execution of the program's functions.
func run(filename string, out io.Writer, skipPreview bool) error {
	// Read data from the input file and check for errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	htmlData := parseContent(input)

	// Create a temp file and check for errors
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()

	fmt.Fprintln(out, outName)

	if err = saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview(outName)
}

func parseContent(input []byte) []byte {
	// Parse the markdown file through blackfriday and
	// bluemonday to generate a valid and safe HTML file
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer

	// Write html to bytes buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outFname string, data []byte) error {
	// Write the bytes to the file.
	return os.WriteFile(outFname, data, 0644)
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("os not supported")
	}

	// Append filename to parameters slice
	cParams = append(cParams, fname)

	// Locate executable in path
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Open the file using default executable
	exec.Command(cPath, cParams...).Run()

	// Give the browser time to open the file before deleting it.
	// This is a temporary solution. We should implement handling
	// and OS signal in the future.
	time.Sleep(2 * time.Second)
	return err
}
