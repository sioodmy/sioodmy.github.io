package main

import (
	"fmt"
	"html/template"
	"os"
	"slices"
	"strings"

	"github.com/sioodmy/website/internal/post"
)

type BlogIndexTemplate struct {
	Posts []post.BlogPostMeta
}

func TitleToFilename(title string) string {
	lower := strings.ToLower(title)
	trimmed := strings.ReplaceAll(lower, " ", "-")
	return fmt.Sprintf("%s.html", trimmed)
}

func GenerateBlog(posts []post.BlogPostRaw) {
	tpl, err := template.New("").ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	var index BlogIndexTemplate
	for _, blogpost := range posts {

		// TODO: parse file name

		filename := TitleToFilename(blogpost.Title)

		out := fmt.Sprintf("website/blog/%s", filename)

		file, _ := os.Create(out)
		defer file.Close()

		date := post.ParseDate(blogpost.Date)
		templatedata := post.BlogPostTemplate{
			DateString: date,
			Title:      blogpost.Title,
			Html:       post.MdToHTML(blogpost.Content),
		}

		tpl.ExecuteTemplate(file, "post.html", templatedata)

		meta := post.BlogPostMeta{
			Date:  date,
			Url:   fmt.Sprintf("blog/%s", filename),
			Title: blogpost.Title,
		}
		index.Posts = append(index.Posts, meta)
	}
	indexfile, _ := os.Create("website/blog.html")
	defer indexfile.Close()
	tpl.ExecuteTemplate(indexfile, "index.html", index)

}
func main() {
	items, _ := os.ReadDir("./blog")

	var blogposts []post.BlogPostRaw
	for _, item := range items {
		path := fmt.Sprintf("blog/%s", item.Name())
		blogpost := post.FromFile(path)
		blogposts = append(blogposts, blogpost)
	}

	compare := func(a, b post.BlogPostRaw) int {
		return a.Date.Compare(b.Date)
	}

	slices.SortFunc(blogposts, compare)

	GenerateBlog(blogposts)
}
