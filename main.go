package main

import (
	"io/fs"
	"log"
	"path/filepath"

	"github.com/acrobatstick/acrobatstick.github.io/markdown"
)

func main() {
	paths := scan("articles")
	for _, p := range paths {
		doc, err := markdown.Read(p)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if err := doc.WriteIntoHTML(); err != nil {
			log.Println(err)
			continue
		}
	}
}

func scan(p string) []string {
	paths := []string{}
	filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(d.Name()) != ".md" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	return paths
}
