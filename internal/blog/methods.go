package blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/2bitburrito/hpa-website/internal/helpers"
)

func NewBlogs() Blogs {
	return make(Blogs)
}

// WriteBlogDataToJSON writes to json file for runtime information
func WriteBlogDataToJSON(blogs Blogs) error {
	dat, err := json.Marshal(blogs)
	if err != nil {
		return err
	}

	err = os.WriteFile(JSONDataFilepath, dat, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func ReadBlogData() (Blogs, error) {
	dat, err := os.ReadFile(JSONDataFilepath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read blog data file: %w", err)
	}

	var bgs Blogs
	err = json.Unmarshal(dat, &bgs)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal blog data file: %w", err)
	}
	return bgs, nil
}

func (b Blogs) Get(name string) (Blog, bool) {
	blog, ok := b[name]
	return blog, ok
}

// New creates a new blog from the data and adds it to the Blogs map
func (b Blogs) New(data map[string]interface{}, fileName string) error {
	name, ok := strings.CutSuffix(fileName, ".md")
	if !ok {
		return fmt.Errorf("couldn't cut %q from file: %s", ".md", fileName)
	}

	title, ok := data["title"].(string)
	if !ok {
		return fmt.Errorf("blog missing required field: title")
	}
	date, ok := data["date"].(string)
	if !ok {
		date = ""
	}
	description, ok := data["description"].(string)
	if !ok {
		description = ""
	}
	isDraft, ok := data["draft"].(bool)
	if !ok {
		isDraft = false
	}

	blg := Blog{
		Title:       title,
		Description: description,
		Date:        date,
		HTMLContent: bytes.Buffer{},
		IsDraft:     isDraft,
		Filepath:    fmt.Sprintf("%s/blog/articles/%s.html", helpers.OutDir, name),
		FileName:    name,
	}

	b[name] = blg

	return nil
}
