package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type item struct {
	Prefix      string `json:"prefix"`
	Description string `json:"description"`
}

type snippet map[string]item

type snippetfile struct {
	Filename string
	snippets snippet
}

func main() {
	// Load file names
	var files []string
	filepath.Walk("./snippets/", func(file string, info os.FileInfo, err error) error {
		if filepath.Ext(file) == ".code-snippets" {
			files = append(files, file)
		}

		return nil
	})

	// Transform in snippet files instances
	var snippets []snippetfile
	for _, file := range files {
		snippets = append(snippets, parse(file))
	}

	// Write catalog file
	dumpToMarkdown(snippets)
}

func parse(path string) snippetfile {
	_, filenameWithExt := filepath.Split(path)
	name := strings.TrimSuffix(filenameWithExt, filepath.Ext(path))

	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var snippet snippet

	json.Unmarshal([]byte(bytes), &snippet)

	return snippetfile{
		Filename: name,
		snippets: snippet,
	}
}

func dumpToMarkdown(files []snippetfile) {
	f, err := os.OpenFile("CATALOG.md", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	f.WriteString("# Catalog\n")
	f.WriteString("\n")
	f.WriteString("Here's a list of snippets available at the moment:\n")
	f.WriteString("\n")

	for i := 0; i < len(files); i++ {
		file := files[i]

		f.WriteString(fmt.Sprintf("## %s\n", file.Filename))
		f.WriteString("\n")

		f.WriteString("| Prefix | Description |\n")
		f.WriteString("|--------|-------------|\n")

		for _, v := range file.snippets {
			f.WriteString(fmt.Sprintf("| %s | %s |\n", v.Prefix, v.Description))
		}

		f.WriteString("\n")
		f.WriteString("---")
		f.WriteString("\n\n")
	}
}
