package blog

import (
	"bytes"
)

const JSONDataFilepath = "./static/blog/data.json"

type Blogs map[string]Blog

type Blog struct {
	Title       string
	Description string
	Date        string
	HTMLContent bytes.Buffer
	IsDraft     bool
	Filepath    string // Path of generated HTML File
	FileName    string // This is the name without ".md". It is also the url path
}
