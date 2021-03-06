package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const defaultTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>{{ .Title }}</title>
  </head>
  <body>
{{ .Body }}
  </body>
</html>
`

// content represents the HTML content to add into the template.
type content struct {
	Title string
	File  string
	Body  template.HTML
}

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	stdin := flag.Bool("stdin", false, "Read from stdin")
	flag.Parse()

	if !*stdin {
		// If no file input provided, show usage.
		if *filename == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

	if err := run(*filename, *tFname, os.Stdout, *skipPreview, *stdin); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run coordinates te execution of the program's functions.
func run(filename, tFname string, out io.Writer, skipPreview, stdin bool) (err error) {
	var input []byte

	if !stdin {
		// Read data from the input file and check for errors
		input, err = os.ReadFile(filename)
		if err != nil {
			return err
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		filename = "stdin"
	}

	htmlData, err := parseContent(input, filename, tFname)
	if err != nil {
		return err
	}

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

func parseContent(input []byte, srcFileName, tFname string) ([]byte, error) {
	// Parse the markdown file through blackfriday and
	// bluemonday to generate a valid and safe HTML file
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Parse content of the defaultTemplate const into a new Template
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	// If user provides an alternate template, replace default.
	// Use can also use env var to set filename.
	if f := os.Getenv("TEMPLATE_FILENAME"); f != "" {
		tFname = os.Getenv("TEMPLATE_FILENAME")
	}
	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	// Instantiate the content type, adding the title and body.
	c := content{
		Title: "Markdown Preview Tool",
		File:  filepath.Base(srcFileName),
		Body:  template.HTML(body),
	}

	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer

	// Execute the template with the content type
	if err = t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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
