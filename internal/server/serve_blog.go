package server

import (
	"fmt"
	"net/http"
)

func (s *Server) HandleServeBlog(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	blog, exists := s.Dependencies.Blogs.Get(name)
	if !exists {
		http.ServeFile(w, r, "./static/404.html")
		return
	}

	if blog.IsDraft && !s.isDev {
		http.ServeFile(w, r, "./static/404.html")
		return
	}
	fmt.Println("serving: ", blog.Filepath)
	http.ServeFile(w, r, blog.Filepath)
}
