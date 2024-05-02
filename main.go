package main

import (
	"os"

	"github.com/sioodmy/generator/internal/index"
	"github.com/sioodmy/generator/internal/projects"
)

func main() {
	os.MkdirAll("generated/blog", 0o777)

	index.GenerateBlog(index.GetPosts())
	projects.GenerateProjects()
}
