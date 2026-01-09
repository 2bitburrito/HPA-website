package main

import (
	"fmt"

	"github.com/2bitburrito/hpa-website/cmd/generate/templator"
	"github.com/2bitburrito/hpa-website/internal/blog"
)

func main() {
	fmt.Println("Building Website...")

	err := makeDirs()
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
