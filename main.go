package main

import (
	"fmt"
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
)

func processTemplate(sourceFile, destinationFile string) error {
	t := template.Must(template.ParseFiles(sourceFile))
	f, err := os.Create(strings.TrimSuffix(destinationFile, templateSuffix))
	if err != nil {
		return err
	}
	err = t.Execute(f, variables)
	if err != nil {
		return err
	}
	return f.Close()
}

func processFile(sourceInfo os.FileInfo, sourceFile, destinationFile string) error {
	if sourceInfo.IsDir() {
		return nil
	}
	destinationFile = strings.ReplaceAll(destinationFile, application, variables.Name)
	for _, ignoreFile := range ignoreList {
		if ignoreFile == sourceFile {
			log.Println("ignore " + ignoreFile)
		}
	}
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}
	fmt.Println("." + destinationFile)
	err = os.MkdirAll(path.Dir(destinationFile), 0755)
	if err != nil {
		return err
	}
	if strings.HasSuffix(sourceFile, templateSuffix) {
		return processTemplate(sourceFile, destinationFile)
	} else {
		return ioutil.WriteFile(destinationFile, input, 0644)
	}
}

func main() {
	targetDirectory := "output"
	fmt.Println(os.Args)
	if len(os.Args) > 1 {
		targetDirectory = os.Args[1]
		if len(os.Args) > 2 {
			variables.Name = os.Args[2]
		}
		if len(os.Args) > 3 {
			targetPlatform = os.Args[3]
		}
	}
	sourceRoot := filepath.Join(".", targetPlatform)
	err := filepath.Walk(sourceRoot,
		func(sourcePath string, sourceInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			targetFile := filepath.Join(targetDirectory, sourcePath)
			return processFile(sourceInfo, sourcePath, targetFile)
		})
	if err != nil {
		log.Println(err)
	}
}
