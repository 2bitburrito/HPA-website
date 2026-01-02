package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/2bitburrito/hpa-website/internal/blog"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
)

type templator struct {
	mdRenderer goldmark.Markdown
	scaffold   HTMLScaffold
	Blogs      blog.Blogs
}

type HTMLScaffold struct {
	Head   string
	Footer string
	NavBar string
}

func newTemplator() (*templator, error) {
	t := templator{
		mdRenderer: goldmark.New(
			goldmark.WithExtensions(
				meta.Meta,
			),
		),
	}

	scaffold, err := t.scaffoldHTML()
	if err != nil {
		return nil, err
	}

	err = t.GenerateBlogFiles()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// scaffoldHTML reads the snippet html for generic pages, writes them to the main and blog-index
// It also returns the scaffold to be saved into the templator and used for future files as needed
func (t *templator) scaffoldHTML() error {
	head, err := os.ReadFile(TemplatesDirectory + "snippets/head.html")
	if err != nil {
		return err
	}
	footer, err := os.ReadFile(TemplatesDirectory + "snippets/foot.html")
	if err != nil {
		return err
	}
	navBar, err := os.ReadFile(TemplatesDirectory + "snippets/nav-bar.html")
	if err != nil {
		return err
	}
	t.scaffold = HTMLScaffold{
		Head:   string(head),
		Footer: string(footer),
		NavBar: string(navBar),
	}

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
