package post

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type BlogPostRaw struct {
	Date    time.Time
	Title   string
	Content string
}

type BlogPostMeta struct {
	Date, Url, Title string
}

type BlogPostTemplate struct {
	Title, DateString string
	Html              template.HTML
}

func FromFile(path string) BlogPostRaw {
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	content := string(f)
	data := strings.SplitN(string(content), "\n", 3)
	if len(data) < 3 {
		panic("Incorrect blogpost syntax")
	}
	sec, err := strconv.ParseInt(strings.TrimSpace(data[0]), 10, 64)
	if err != nil {
		panic(err)
	}
	blogpost := BlogPostRaw{

		Date:    time.Unix(sec, 0),
		Title:   data[1][2:],
		Content: data[2],
	}
	return blogpost
}

func ParseDate(date time.Time) string {
	return fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())
}

func MdToHTML(md string) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return template.HTML(string(markdown.Render(doc, renderer)))
}
