package main

import (
	"fmt"
	"os"

	"github.com/2bitburrito/hpa-website/cmd/generate/templator"
	"github.com/2bitburrito/hpa-website/internal/blog"
	"github.com/2bitburrito/hpa-website/internal/helpers"
)

func main() {
	fmt.Println("Building Website...")

	err := os.MkdirAll(helpers.OutDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	t, err := templator.New()
	if err != nil {
		panic(err)
	}

	err = t.WriteHTMLArticles()
	if err != nil {
		panic(err)
	}

	err = blog.WriteBlogDataToJSON(t.Blogs)
	if err != nil {
		panic(err)
	}

	fmt.Println("Website Generated Successfully")
}
