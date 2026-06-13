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

	fmt.Printf("\n\033[32midkspg - I Don't Know Static Page Generator \n ———————————————————————————————— \n\033[0m\n")
	
	// Reads template file
	fmt.Printf("Reading templates... \n\n")
	templateHtml, err := os.ReadFile("template.html")
	itemHtml, err := os.ReadFile("item.html")
	indexHtml, err := os.ReadFile("indextemplate.html")
	
	// Continues operation only if template.html is read
	if(err != nil){
		fmt.Printf("\n\033[31mError reading templates \n%v\n", err)
	}else{
		traverse("blog/", templateHtml, itemHtml, indexHtml)
	}
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

func traverse(rootpath string, templatehtml []byte, itemHtml []byte, indexHtml []byte){

	var indexItems string

	// Traverses the directory
	err := filepath.WalkDir(rootpath, func(path string, info fs.DirEntry, err error) error {

		if filepath.Ext(path) == ".md" {
			// Prints the current working file
			fmt.Printf("Processing %s\n", path)

			// Opens the .md file
			mdfile, err := os.ReadFile(path)
			if(err != nil){
				fmt.Printf("\n\033[31mError reading .md file \n%v\n", err)
			}else{
				// Joins the parsed file with a template.html 
				html := bytes.Replace(templatehtml, []byte("<!-- REPLACE -->"), mdToHTML(mdfile), -1)

				// Writes the file to the working directory
				outfile := filepath.Dir(path) + "/index.html"
				os.WriteFile(outfile, html, 0644)

				indexItems = indexItems + string(itemBuilder(scanMetadata(path), outfile, itemHtml))
			}
		}
		return nil
	})

	html := bytes.Replace(indexHtml, []byte("<!-- REPLACE -->"), []byte(indexItems), -1)
	os.WriteFile("index.html", html, 0644)

	if err != nil {
		fmt.Printf("womp womp: %v\n", err)
	}
}

func scanMetadata(path string) map[string]string{

	metadata := make(map[string]string)
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

func itemBuilder(metadata map[string]string, link string, itemHtml []byte) []byte {


	item := bytes.Replace(itemHtml, []byte("<!-- TITLE -->"), []byte(metadata["_TITLE"]), -1)
	item = bytes.Replace(item, []byte("<!-- DESCRIPTION -->"), []byte(metadata["_DESCRIPTION"]), -1)
	item = bytes.Replace(item, []byte("<!-- LINK -->"), []byte(link), -1)
	return item
}