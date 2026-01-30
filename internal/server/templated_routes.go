package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/2bitburrito/hpa-website/internal/blog"
	"github.com/2bitburrito/hpa-website/internal/helpers"
)

type RenderData struct {
	blog.Blog
	Views int
}

func (s *Server) HandleMainPage(w http.ResponseWriter, r *http.Request) {
	// Increment the view count for the home page
	s.Dependencies.SheetsService.IncrementMain()

	lstTmp, err := template.ParseFiles(helpers.TemplatesDirectory + "snippets/blog-list.html")
	if err != nil {
		fmt.Println("failed to parse main.html template: ", err)
		http.ServeFile(w, r, helpers.OutDir+"404.html")
		return
	}

	blogs := s.Dependencies.Blogs.GetNum(5, s.isDev)

	// Add the views to the blogs data
	dat, err := s.addViewsToBlogs(blogs)
	if err != nil {
		fmt.Println("failed to add views to blogs: ", err)
		http.ServeFile(w, r, helpers.OutDir+"404.html")
		return
	}

	var buf bytes.Buffer
	err = lstTmp.Execute(&buf, dat)
	if err != nil {
		fmt.Println("failed to execute main.html template: ", err)
		http.ServeFile(w, r, helpers.OutDir+"404.html")
		return
	}
	wholeLst := template.HTML(buf.String())

	mainTmp := template.New("index.html").Delims("[[", "]]")
	mainTmp, err = mainTmp.ParseFiles(helpers.OutDir + "main/index.html")
	if err != nil {
		fmt.Println("failed to parse main.html template: ", err)
		http.ServeFile(w, r, helpers.OutDir+"404.html")
		return
	}

	err = mainTmp.Execute(w, wholeLst)
	if err != nil {
		fmt.Println("failed to execute main.html template: ", err)
		http.ServeFile(w, r, helpers.OutDir+"404.html")
		return
	}
}

func (s *Server) HandleBlogIndex(w http.ResponseWriter, r *http.Request) {
	// First render the blog list
	lstTmp, err := template.ParseFiles(helpers.TemplatesDirectory + "snippets/blog-list.html")
	if err != nil {
		fmt.Println("failed to parse blog-index.html template: ", err)
		http.ServeFile(w, r, helpers.OutDir+"404.html")
		return
	}

	blogs := s.Dependencies.Blogs.GetNum(30, s.isDev)

	// Add the views to the blogs data
	dat, err := s.addViewsToBlogs(blogs)
	if err != nil {
		fmt.Println("failed to add views to blogs: ", err)
		http.ServeFile(w, r, helpers.OutDir+"404.html")
		return
	}

	var buf bytes.Buffer
	err = lstTmp.Execute(&buf, dat)
	if err != nil {
		fmt.Println("failed to execute blog-index.html template: ", err)
		http.ServeFile(w, r, helpers.OutDir+"404.html")
		return
	}
	wholeLst := template.HTML(buf.String())

	mainTmp := template.New("index.html").Delims("[[", "]]")
	mainTmp, err = mainTmp.ParseFiles(helpers.OutDir + "blog/index.html")
	if err != nil {
		fmt.Println("failed to parse blog-index.html template: ", err)
		http.ServeFile(w, r, helpers.OutDir+"404.html")
		return
	}

	err = mainTmp.Execute(w, wholeLst)
	if err != nil {
		fmt.Println("failed to execute blog-index.html template: ", err)
		http.ServeFile(w, r, helpers.OutDir+"404.html")
		return
	}
}

func (s *Server) addViewsToBlogs(blogs []blog.Blog) ([]RenderData, error) {
	dat := make([]RenderData, 0, len(blogs))
	for _, b := range blogs {
		views, err := s.Dependencies.SheetsService.GetViews(b.FileName)
		if err != nil {
			return nil, fmt.Errorf("failed to get views for blog %s: %w", b.FileName, err)
		}
		dat = append(dat, RenderData{
			Blog:  b,
			Views: views,
		})
	}
	return dat, nil
}
