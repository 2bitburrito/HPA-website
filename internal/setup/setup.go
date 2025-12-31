// Package setup is the dependency struct
package setup

import (
	"context"

	"github.com/2bitburrito/hpa-website/internal/blog"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Dependencies struct {
	Aws   aws.Config
	Blogs blog.Blogs
}

func Setup() (Dependencies, error) {
	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("ap-southeast-2"),
	)
	if err != nil {
		return Dependencies{}, err
	}

	blogs, err := blog.ReadBlogData()
	if err != nil {
		return Dependencies{}, err
	}

	return Dependencies{
		Aws:   awsCfg,
		Blogs: blogs,
	}, nil
}
