package markdown

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/acrobatstick/acrobatstick.github.io/components"
	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gosimple/slug"
)

type Document struct {
	Title         string    `yaml:"title"`
	DatePublished time.Time `yaml:"date"`
	Slug          string    `yaml:"slug"`
	RawMarkdown   []byte
	htmlBuf       bytes.Buffer
	path          string
}

func Read(p string) (*Document, error) {
	doc := new(Document)
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)
	rest, err := frontmatter.Parse(r, &doc)
	if err != nil {
		return nil, err
	}
	if len(doc.Slug) == 0 {
		basename := strings.TrimSuffix(path.Base(p), path.Ext(p))
		s := slug.Make(basename)
		doc.Slug = s
	}
	doc.RawMarkdown = rest
	doc.path = p
	return doc, nil
}

func (d *Document) WriteIntoHTML() error {
	basename := strings.TrimSuffix(path.Base(d.path), path.Ext(d.path))
	content := string(parseToHTML(d.RawMarkdown))
	var buf bytes.Buffer
	err := components.Layout(d.Title, toTempl(content)).Render(context.Background(), &buf)
	if err != nil {
		return err
	}

	dir := path.Join("dist", "article", d.DatePublished.Format("2006/01"))
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("mkdir %q failed for article %q: %w", dir, basename, err)
	}

	outPath := path.Join(dir, d.Slug+".html")
	err = os.WriteFile(outPath, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("writefile to %q: %w", outPath, err)
	}

	return nil
}

func parseToHTML(md []byte) []byte {
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

func toTempl(content string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, content)
		return
	})
}
