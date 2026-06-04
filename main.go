package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"fmt"
	"os"
	"path/filepath"
	"io/fs"
	"bytes"
	"bufio"
	"strings"
)

func main() {

	templatehtml, _ := os.ReadFile("template.html")
	traverse("blog/", templatehtml)

}

func mdToHTML(md []byte) []byte {
	
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func traverse(rootpath string, templatehtml []byte){

	// Initialize metadata variable for the .md files
	metadata := make(map[string]string)

	// Traverses the directory
	err := filepath.WalkDir(rootpath, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() && path != rootpath {
			// Prints the current working directory
			fmt.Printf("Processing %s\n", path)
		}

		if filepath.Ext(path) == ".md" {

			// Opens the .md file
			mdfile, _ := os.ReadFile(path)
			// Joins the parsed file with a template.html 
			html := bytes.Replace(templatehtml, []byte("<!-- REPLACE -->"), mdToHTML(mdfile) , -1)
			// Writes the file to the working directory
			os.WriteFile(filepath.Dir(path) + "/index.html", html, 0644)
		
		}
		return nil
	})

	if err != nil {
		fmt.Printf("womp womp: %v\n", err)
	}
}

func scanMetadata(path []byte) map[string]string{

	file, _ := os.Open(path)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Filter lines that start with _ prefix
		if strings.HasPrefix(line, "_") {

			parts := strings.SplitN(line, "=", 2)

			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
			
				metadata[key] = value
			}
		}
	}
	return metadata
}