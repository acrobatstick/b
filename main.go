package main

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/acrobatstick/acrobatstick.github.io/components"
	"github.com/acrobatstick/acrobatstick.github.io/markdown"
	"github.com/gosimple/slug"
	"github.com/yosssi/gohtml"
)

var (
	ROOT_PATH = "public"
)

func main() {
	paths := scanArticles("articles")
	for _, p := range paths {
		b, err := markdown.Read(p)
		if err != nil {
			panic(err)
		}
		if err := writeArticle(p, b); err != nil {
			log.Printf("Error writing article: %v", err)
		}
	}
}

func scanArticles(p string) []string {
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

func writeArticle(p string, b []byte) error {
	basename := strings.TrimSuffix(path.Base(p), path.Ext(p))
	s := slug.Make(basename)
	content := string(markdown.RenderToHTML(b))

	var buf bytes.Buffer
	err := components.ArticlePage(basename, toTempl(content)).Render(context.Background(), &buf)
	if err != nil {
		return err
	}

	formatted := gohtml.FormatBytes(buf.Bytes())
	outPath := path.Join(ROOT_PATH, s+".html")
	err = os.WriteFile(outPath, formatted, 0644)
	if err != nil {
		return err
	}

	return nil
}

func toTempl(content string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, content)
		return
	})
}
