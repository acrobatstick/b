package main

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/acrobatstick/acrobatstick.github.io/components"
	"github.com/acrobatstick/acrobatstick.github.io/markdown"
)

func main() {
	if err := build(); err != nil {
		panic(err)
	}
}

func build() error {
	if _, err := os.Stat("dist"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("dist", 077)
		}
	}
	name := path.Join("dist", "index.html")
	// build the index page
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	err = components.Index().Render(context.Background(), f)
	if err != nil {
		return err
	}
	// scan for articles
	paths := scan("articles")
	for _, p := range paths {
		doc, err := markdown.Read(p)
		if err != nil {
			log.Println(err)
			continue
		}
		if err := doc.WriteIntoHTML(); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
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
