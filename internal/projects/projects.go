package projects

import (
	"html/template"
	"os"

	"github.com/pelletier/go-toml"
)

type Project struct {
	Name, Description, RepoUrl, BlogpostUrl string
}
type ProjectsPage struct {
	Projects []Project
}

func GenerateProjects() {
	tpl, _ := template.New("").ParseGlob("templates/projects.html")

	f, _ := os.Create("generated/projects.html")

	config, _ := os.ReadFile("projects.toml")

	var data ProjectsPage
	err := toml.Unmarshal(config, &data)
	if err != nil {
		panic(err)
	}
	tpl.ExecuteTemplate(f, "projects.html", &data)
}
