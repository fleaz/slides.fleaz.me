package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/gographics/imagick.v2/imagick"
)

type Talk struct {
	Title       string
	Preview     string
	Path        string
	Description string
}

func generatePreview(path string) Talk {
	r := regexp.MustCompile("talks\\/(.*)\\.pdf")
	result := r.FindStringSubmatch(path)
	imagePath := fmt.Sprintf("previews/%s.jpg", result[1])
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err := mw.ReadImage(path)
	if err != nil {
		panic(err)
	}
	mw.SetIteratorIndex(0)
	mw.SetImageFormat("jpg")
	err = mw.WriteImage(imagePath)
	if err != nil {
		panic(err)
	}
	return Talk{Title: result[1], Preview: imagePath, Path: path, Description: "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."}
}

func main() {

	var files []string

	root := "talks/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	imagick.Initialize()
	defer imagick.Terminate()
	var data []Talk
	for _, file := range files {
		t := generatePreview(file)
		data = append(data, t)
	}

	for i := len(data)/2 - 1; i >= 0; i-- {
		opp := len(data) - 1 - i
		data[i], data[opp] = data[opp], data[i]
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		panic(err)
	}

	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	tmpl.Execute(f, data)
}
