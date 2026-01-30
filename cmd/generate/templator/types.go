package templator

import (
	"html/template"

	"github.com/2bitburrito/hpa-website/internal/blog"
	"github.com/yuin/goldmark"
)

type RenderData struct {
	HTMLScaffold
	blog.Blog
}

type templator struct {
	isDev      bool
	mdRenderer goldmark.Markdown
	scaffold   HTMLScaffold
	Blogs      blog.Blogs
}

type HTMLScaffold struct {
	Head            template.HTML
	Footer          template.HTML
	NavBar          template.HTML
	BlogListPreview template.HTML
	BlogList        template.HTML
}
