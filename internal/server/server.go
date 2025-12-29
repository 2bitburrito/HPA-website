package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/2bitburrito/hpa-website/internal/setup"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port         int
	Dependencies setup.Dependencies
	Blogs        []Blog
}

func NewServer(params setup.Dependencies) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:         port,
		Dependencies: params,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Println("Server Started Successfully On Port: ", port)
	return server
}
