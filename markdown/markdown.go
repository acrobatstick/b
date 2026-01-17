package markdown

import (
	"os"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Document struct {
	Title   string
	RawHTML []byte
	Date    time.Time
}

func Read(p string) ([]byte, error) {
	// TODO: get the file creation date to order the blog list later
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func RenderToHTML(md []byte) []byte {
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
