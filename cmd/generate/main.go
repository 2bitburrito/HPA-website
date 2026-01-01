package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/2bitburrito/hpa-website/internal/blog"
	"github.com/2bitburrito/hpa-website/internal/helpers"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

const (
	TemplatesDirectory = "./templates/"
	BlogPath           = "./static/blog-md/"
)

type templator struct {
	mdRenderer goldmark.Markdown
}

func main() {
	fmt.Println("Building Website...")
	t := templator{
		mdRenderer: goldmark.New(
			goldmark.WithExtensions(
				meta.Meta,
			),
		),
	}

	fls, err := os.ReadDir(BlogPath)
	if err != nil {
		panic(err)
	}

	blogs, err := t.GenerateBlogFiles(fls)
	if err != nil {
		panic(err)
	}

	err = formatBlogIndexSnippet(blogs)
	if err != nil {
		panic(err)
	}

	err = writeHTMLArticleFiles(blogs)
	if err != nil {
		panic(err)
	}

	err = blog.WriteBlogDataToJSON(blogs)
	if err != nil {
		panic(err)
	}
	fmt.Println("Website Generated Successfully")
}

func (t *templator) GenerateBlogFiles(files []os.DirEntry) (blog.Blogs, error) {
	blogs := blog.NewBlogs()

	var buf bytes.Buffer
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		mdn, err := os.ReadFile(BlogPath + f.Name())
		if err != nil {
			return blog.Blogs{}, fmt.Errorf("failed to read file: %s\n%v", f.Name(), err)
		}

		buf.Reset()

		ctx := parser.NewContext()
		err = t.mdRenderer.Convert(mdn, &buf, parser.WithContext(ctx))
		if err != nil {
			return blog.Blogs{}, fmt.Errorf("failed to convert file: %s\n%v", f.Name(), err)
		}

		data := meta.Get(ctx)
		err = blogs.AddNew(data, f.Name(), buf)
		if err != nil {
			return blog.Blogs{}, fmt.Errorf("failed to create blog: %s\n%v", f.Name(), err)
		}
	}
	return *blogs, nil
}

func writeHTMLArticleFiles(blogs blog.Blogs) error {
	tmpl, err := template.ParseFiles(TemplatesDirectory + "article.gohtml")
	if err != nil {
		return fmt.Errorf("failed to parse article template %w", err)
	}

	for _, b := range blogs {
		err = renderAndSave(tmpl, b, b.Filepath)
		if err != nil {
			return err
		}
	}
	return nil
}

// Formats the main file and the blog index with the data
func formatBlogIndexSnippet(blogs blog.Blogs) error {
	tmpl, err := template.ParseFiles(TemplatesDirectory + "snippets/blog_list.html")
	if err != nil {
		return fmt.Errorf("failed to parse blog_list.html template: %w", err)
	}

	// Render the first 5 articles for main page
	err = renderAndSave(tmpl, blogs.Limit(5), helpers.OutDir+"lib/blog_list_preview.html")
	if err != nil {
		return err
	}

	// Render all of the articles for the blog index
	err = renderAndSave(tmpl, blogs, helpers.OutDir+"lib/blog_list.html")
	if err != nil {
		return err
	}
	return nil
}

func renderAndSave(tmpl *template.Template, data any, filepath string) error {
	var buf bytes.Buffer

	err := tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	err = os.WriteFile(filepath, buf.Bytes(), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
