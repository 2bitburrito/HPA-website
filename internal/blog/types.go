package blog

import (
	"html/template"
	"time"
)

const JSONDataFilepath = "./static/blog/data.json"

type Blogs []Blog

type Blog struct {
	Title       string
	Description string
	Date        time.Time
	HTMLContent template.HTML
	IsDraft     bool
	Filepath    string // Path of generated HTML File
	FileName    string // This is the name without ".md". It is also the url path
}
