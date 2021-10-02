package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Variables struct {
	Name   string
	Target string
}

const (
	templateSuffix = ".tmpl"
	application    = "grandana"
)

var (
	dryrun         = true
	targetPlatform = "gcp"
	variables      = Variables{
		Name: "grandana",
	}
	ignoreList = []string{
		"main",
		"main.go",
		".gitignore",
	}
	sourceRoot = filepath.Join(".", targetPlatform)
)

func processTemplate(sourceFile, destinationFile string) error {
	t := template.Must(template.ParseFiles(sourceFile))
	f, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	err = t.Execute(f, variables)
	if err != nil {
		return err
	}
	log.Println(destinationFile)
	return f.Close()
}

func processFile(sourceInfo os.FileInfo, sourceFile, destinationFile string) error {
	if sourceInfo.IsDir() {
		return nil
	}
	destinationFile = strings.ReplaceAll(destinationFile, application, variables.Name)
	for _, ignoreFile := range ignoreList {
		if strings.Contains(sourceFile, ignoreFile) {
			log.Println("ignoring " + ignoreFile)
			return nil
		}
	}
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Dir(destinationFile), 0755)
	if err != nil {
		return err
	}
	if strings.HasSuffix(sourceFile, templateSuffix) {
		return processTemplate(sourceFile, strings.TrimSuffix(destinationFile, templateSuffix))
	}
	log.Println(destinationFile)
	return ioutil.WriteFile(destinationFile, input, 0644)
}

func main() {
	targetDirectory := "output"
	if len(os.Args) > 1 {
		targetDirectory = os.Args[1]
		if len(os.Args) > 2 {
			variables.Name = os.Args[2]
		}
		if len(os.Args) > 3 {
			targetPlatform = os.Args[3]
		}
	}
	err := filepath.Walk(sourceRoot,
		func(sourcePath string, sourceInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			destinationFile := filepath.Join(targetDirectory, strings.TrimPrefix(sourcePath, sourceRoot))
			return processFile(sourceInfo, sourcePath, destinationFile)
		})
	if err != nil {
		log.Println(err)
	}
}
