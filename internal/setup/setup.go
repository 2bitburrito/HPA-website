// Package setup is the dependency struct
package setup

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Dependencies struct {
	Aws aws.Config
}

func Setup() (Dependencies, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("ap-southeast-2"))
	if err != nil {
		return Dependencies{}, err
	}
	return Dependencies{Aws: awsCfg}, nil
}
