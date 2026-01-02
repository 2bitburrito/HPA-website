package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/2bitburrito/hpa-website/internal/blog"
	"github.com/2bitburrito/hpa-website/internal/helpers"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

const (
	TemplatesDirectory = "./templates/"
	BlogPath           = "./static/blog-md/"
)

type RenderData struct {
	HTMLScaffold HTMLScaffold
	Blog         blog.Blog
}

func main() {
	fmt.Println("Building Website...")

	t, err := newTemplator()
	if err != nil {
		panic(err)
	}

	err = scaffoldHTML()
	if err != nil {
		panic(err)
	}

	err := t.GenerateBlogFiles()
	if err != nil {
		panic(err)
	}

	err = t.writeHTMLArticleFiles()
	if err != nil {
		panic(err)
	}

	err = blog.WriteBlogDataToJSON(t.Blogs)
	if err != nil {
		panic(err)
	}

	err = t.formatBlogIndexSnippets()
	if err != nil {
		panic(err)
	}

	fmt.Println("Website Generated Successfully")
}

func (t *templator) GenerateBlogFiles() error {
	files, err := os.ReadDir(BlogPath)
	if err != nil {
		return fmt.Errorf("failed to read blog directory: %w", err)
	}

	blogs := blog.NewBlogs()

	var buf bytes.Buffer
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		mdn, err := os.ReadFile(BlogPath + f.Name())
		if err != nil {
			return fmt.Errorf("failed to read file: %s\n%v", f.Name(), err)
		}

		buf.Reset()

		ctx := parser.NewContext()
		err = t.mdRenderer.Convert(mdn, &buf, parser.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("failed to convert file: %s\n%v", f.Name(), err)
		}

		data := meta.Get(ctx)
		err = blogs.AddNew(data, f.Name(), buf)
		if err != nil {
			return fmt.Errorf("failed to create blog: %s\n%v", f.Name(), err)
		}
	}
	t.Blogs = *blogs
	return nil
}

func (t *templator) writeHTMLArticleFiles() error {
	tmpl, err := template.ParseFiles(TemplatesDirectory + "article.html")
	if err != nil {
		return fmt.Errorf("failed to parse article template %w", err)
	}

	for _, b := range t.Blogs {
		r := RenderData{
			HTMLScaffold: t.scaffold,
			Blog:         b,
		}
		err = renderAndSave(tmpl, r, b.Filepath)
		if err != nil {
			return err
		}
	}
	return nil
}

// Formats the main file and the blog index with the data
func (t *templator) formatBlogIndexSnippets() error {
	tmpl, err := template.ParseFiles(TemplatesDirectory + "snippets/blog_list.html")
	if err != nil {
		return fmt.Errorf("failed to parse blog_list.html template: %w", err)
	}

	// Render the first 5 articles for main page
	err = renderAndSave(tmpl, t.Blogs.Limit(5), helpers.OutDir+"lib/blog_list_preview.html")
	if err != nil {
		return err
	}

	// Render all of the articles for the blog index
	err = renderAndSave(tmpl, t.Blogs, helpers.OutDir+"lib/blog_list.html")
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

func scaffoldHTML() error {
	tmpl, err := template.ParseFiles(
		TemplatesDirectory+"main.html",
		TemplatesDirectory+"blog-index.html",
		TemplatesDirectory+"article.html",
		// TemplatesDirectory+"snippets/head.html",
		// TemplatesDirectory+"snippets/foot.html",
		// TemplatesDirectory+"snippets/nav_bar.html",
	)
	if err != nil {
		return fmt.Errorf("failed to parse main|snippet templates: %w", err)
	}

	return nil
}
