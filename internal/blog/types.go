package blog

import (
	"html/template"
	"time"
)

const JSONDataFilepath = "./static/blog/data.json"

type Blogs []Blog

type Blog struct {
	BaseBlog
	HTMLContent template.HTML
}

// BaseBlog is the same as Blog but without the HTMLContent
// It is made so that we can save to json without all the content
type BaseBlog struct {
	Title       string
	Description string
	Date        time.Time
	IsDraft     bool
	Filepath    string // Path of generated HTML File
	FileName    string // This is the name without ".md". It is also the url path
}
