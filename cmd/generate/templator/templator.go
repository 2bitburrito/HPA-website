package templator

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/2bitburrito/hpa-website/internal/blog"
	"github.com/2bitburrito/hpa-website/internal/helpers"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func New() (*templator, error) {
	t := templator{
		mdRenderer: goldmark.New(
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
			goldmark.WithExtensions(
				meta.Meta,
				extension.GFM,
				highlighting.NewHighlighting(
					highlighting.WithStyle("dracula"),
				),
			),
		),
	}

	err := t.GenerateBlogFiles()
	if err != nil {
		return nil, err
	}

	err = t.formatBlogIndexSnippets()
	if err != nil {
		return nil, err
	}

	err = t.scaffoldHTML()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// scaffoldHTML reads the snippet html for generic pages, writes them to the main and blog-index
// It also returns the scaffold to be saved into the templator and used for future files as needed
func (t *templator) scaffoldHTML() error {
	head, err := os.ReadFile(SnippetDirectory + "head.html")
	if err != nil {
		return err
	}
	footer, err := os.ReadFile(SnippetDirectory + "footer.html")
	if err != nil {
		return err
	}
	navBar, err := os.ReadFile(SnippetDirectory + "nav-bar.html")
	if err != nil {
		return err
	}

	t.scaffold = HTMLScaffold{
		Head:            template.HTML(head),
		Footer:          template.HTML(footer),
		NavBar:          template.HTML(navBar),
		BlogListPreview: t.scaffold.BlogListPreview,
		BlogList:        t.scaffold.BlogList,
	}

	tmpl, err := template.ParseFiles(
		TemplatesDirectory+"main.html",
		TemplatesDirectory+"blog-index.html",
	)
	if err != nil {
		return fmt.Errorf("failed to parse main|snippet templates: %w", err)
	}

	err = renderAndSave(tmpl, t.scaffold, helpers.OutDir+"main/index.html", "main.html")
	if err != nil {
		return fmt.Errorf("failed to render and save index.html: %w", err)
	}

	err = renderAndSave(tmpl, t.scaffold, helpers.OutDir+"blog/index.html", "blog-index.html")
	if err != nil {
		return fmt.Errorf("failed to render and save blog-index.html: %w", err)
	}

	return nil
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

func (t *templator) WriteHTMLArticles() error {
	tmpl, err := template.ParseFiles(TemplatesDirectory + "article.html")
	if err != nil {
		return fmt.Errorf("failed to parse article template %w", err)
	}

	for _, b := range t.Blogs {
		r := RenderData{
			HTMLScaffold: t.scaffold,
			Blog:         b,
		}
		err = renderAndSave(tmpl, r, b.Filepath, "article.html")
		if err != nil {
			return err
		}
	}
	return nil
}

// Formats the main file and the blog index with the data and saves it to the scaffold
func (t *templator) formatBlogIndexSnippets() error {
	tmpl, err := template.ParseFiles(TemplatesDirectory + "snippets/blog-list.html")
	if err != nil {
		return fmt.Errorf("failed to parse blog-list.html template: %w", err)
	}

	var buf bytes.Buffer

	// Render the first 5 articles for main page:
	err = tmpl.Execute(&buf, t.Blogs.Limit(5))
	if err != nil {
		return fmt.Errorf("failed to execute template blog-list.html: %w", err)
	}

	t.scaffold.BlogListPreview = template.HTML(buf.String())

	buf.Reset()

	// Render all of the articles for the blog index:
	err = tmpl.Execute(&buf, t.Blogs)
	if err != nil {
		return fmt.Errorf("failed to execute template blog-list.html: %w", err)
	}
	t.scaffold.BlogList = template.HTML(buf.String())

	return nil
}

// renderAndSave is an internal helper which executes the template and saves the file to filepath
func renderAndSave(tmpl *template.Template, data any, filepath string, name string) error {
	var buf bytes.Buffer

	err := tmpl.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	err = os.WriteFile(filepath, buf.Bytes(), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
