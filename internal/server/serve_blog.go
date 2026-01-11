package server

import (
	"net/http"
)

func (s *Server) HandleServeBlog(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	s.Dependencies.SheetsService.ArticleViews.Increment(name)

	blog, exists := s.Dependencies.Blogs.Get(name, s.isDev)
	if !exists {
		http.ServeFile(w, r, "./static/404.html")
		return
	}

	if blog.IsDraft && !s.isDev {
		http.ServeFile(w, r, "./static/404.html")
		return
	}
	http.ServeFile(w, r, blog.Filepath)
}
