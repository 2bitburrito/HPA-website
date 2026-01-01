package blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/2bitburrito/hpa-website/internal/helpers"
)

func NewBlogs() *Blogs {
	b := make(Blogs, 0)
	return &b
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
	for _, blg := range b {
		if blg.FileName == name {
			return blg, true
		}
	}
	return Blog{}, false
}

// AddNew creates a new blog from the data and adds it to the Blogs slice
func (b *Blogs) AddNew(data map[string]any, fileName string, htmlContent bytes.Buffer) error {
	name, ok := strings.CutSuffix(fileName, ".md")
	if !ok {
		return fmt.Errorf("couldn't cut %q from file: %s", ".md", fileName)
	}

	title, ok := data["title"].(string)
	if !ok {
		return fmt.Errorf("blog missing required field: title")
	}
	dateStr, ok := data["date"].(string)
	if !ok {
		return fmt.Errorf("blog missing required field: date")
	}
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return fmt.Errorf("blog %s couldn't parse date: %w ", name, err)
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
		HTMLContent: template.HTML(htmlContent.String()),
		IsDraft:     isDraft,
		Filepath:    fmt.Sprintf("%sblog/articles/%s.html", helpers.OutDir, name),
		FileName:    name,
	}

	*b = append(*b, blg)
	b.sort()

	return nil
}

func (b Blogs) sort() {
	slices.SortFunc(b, func(a, b Blog) int {
		if a.Date.After(b.Date) {
			return -1
		}
		return 1
	})
}

func (b Blogs) Limit(n int) Blogs {
	if n > len(b) {
		return b
	}
	return b[:n]
}
