package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

const (
	TemplatesDirectory = "./templates/"
	BlogPath           = "./templates/blogs/"
	OutDir             = "./static/generated-files/"
)

type templator struct {
	mdRenderer goldmark.Markdown
}

func main() {
	fmt.Println("Building Website...")
	t := templator{
		mdRenderer: goldmark.New(),
	}

	bgs, err := os.ReadDir(BlogPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("All Blogs", bgs)
	t.GenerateBlogFiles(bgs)
}

func (t *templator) GenerateBlogFiles(bgs []os.DirEntry) {
	for _, b := range bgs {
		mdn, err := os.ReadFile(BlogPath + b.Name())
		if err != nil {
			fmt.Printf("failed to read file: %s\n%v\n", b.Name(), err)
			continue
		}

		var buf bytes.Buffer
		t.mdRenderer.Convert(mdn, &buf)

		n, ok := strings.CutSuffix(b.Name(), ".md")
		if !ok {
			fmt.Printf("couldn't cut %q from file: %s\n", ".md", b.Name())
			continue
		}
		err = writeHTML(buf, n)
		if err != nil {
			fmt.Printf("error while writing html: %v\n", err)
		}

	}
}

func writeHTML(b bytes.Buffer, n string) error {
	newPath := fmt.Sprintf("%s%s.html", OutDir, n)
	_, err := os.Stat(newPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("couldn't check path for existing html file: %w", err)
	}
	err = os.WriteFile(newPath, b.Bytes(), 0o777)
	if err != nil {
		return err
	}
	return nil
}
