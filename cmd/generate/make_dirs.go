package main

import (
	"os"

	"github.com/2bitburrito/hpa-website/internal/helpers"
)

func makeDirs() error {
	err := os.MkdirAll(helpers.OutDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(helpers.OutDir+"main/", os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
