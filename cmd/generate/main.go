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

	blogs := t.GenerateBlogFiles(fls)

	err = formatMainHTML(blogs)
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
}

func (t *templator) GenerateBlogFiles(files []os.DirEntry) blog.Blogs {
	blogs := blog.NewBlogs()

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		mdn, err := os.ReadFile(BlogPath + f.Name())
		if err != nil {
			fmt.Printf("failed to read file: %s\n%v\n", f.Name(), err)
			continue
		}

		var buf bytes.Buffer

		ctx := parser.NewContext()
		err = t.mdRenderer.Convert(mdn, &buf, parser.WithContext(ctx))
		if err != nil {
			fmt.Printf("failed to convert file: %s\n%v\n", f.Name(), err)
		}

		data := meta.Get(ctx)
		err = blogs.New(data, f.Name())
		if err != nil {
			fmt.Printf("failed to create blog: %s\n%v\n", f.Name(), err)
		}
	}
	return blogs
}

func writeHTMLArticleFiles(blogs blog.Blogs) error {
	tmpB, err := os.ReadFile("./templates/article.gohtml")
	if err != nil {
		return fmt.Errorf("failed to read article.gohtml file: %w", err)
	}
	tmpl, err := template.New("article").Parse(string(tmpB))
	if err != nil {
		return fmt.Errorf("failed to parse main template %w", err)
	}
	for _, b := range blogs {
		// First format them with the template
		var fileBuf bytes.Buffer
		err = tmpl.Execute(&fileBuf, b)
		if err != nil {
			return fmt.Errorf("failed to execute template for %s: %w", b.FileName, err)
		}

		// Then write them to the file
		err := os.WriteFile(b.Filepath, b.HTMLContent.Bytes(), 0o644)
		if err != nil {
			return err
		}
	}
	return nil
}

// Formats the main file and the blog index with the data
func formatMainHTML(blogs blog.Blogs) error {
	mainFle, err := os.ReadFile("./templates/main.gohtml")
	if err != nil {
		return fmt.Errorf("failed to read main.gohtml file: %w", err)
	}

	mainTmpl, err := template.New("main").Parse(string(mainFle))
	if err != nil {
		return fmt.Errorf("failed to parse main template %w", err)
	}

	var fileBuf bytes.Buffer

	// Create, and save the main page
	err = mainTmpl.Execute(&fileBuf, blogs)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	err = os.WriteFile(helpers.OutDir+"/main/index.html", fileBuf.Bytes(), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}
	idxFle, err := os.ReadFile("./templates/blog-index.gohtml")
	if err != nil {
		return fmt.Errorf("failed to read main.gohtml file: %w", err)
	}

	// Create and save the blog index
	// TODO: Figure out how to just embed the index into the main page and not have to redo this twice...
	idxTmpl, err := template.New("blog-index").Parse(string(idxFle))
	if err != nil {
		return fmt.Errorf("failed to parse main template %w", err)
	}
	err = idxTmpl.Execute(&fileBuf, blogs)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	err = os.WriteFile(helpers.OutDir+"/blog/index.html", fileBuf.Bytes(), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	return nil
}
