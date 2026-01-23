package server

import (
	"fmt"
	"net/http"
)

func (s *Server) HandleServeBlog(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	err := s.Dependencies.SheetsService.Increment(name)
	if err != nil {
		fmt.Println("failed to increment sheets view count: ", err)
	}

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
