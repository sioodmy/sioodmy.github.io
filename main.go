package main

import (
	"fmt"
	"os"

	"github.com/sioodmy/website/internal/index"
	"github.com/sioodmy/website/internal/projects"
)

func main() {
	os.MkdirAll("generated/blog", 0777)
	fmt.Println(index.GetPosts())

	index.GenerateBlog(index.GetPosts())
	projects.GenerateProjects()
}
