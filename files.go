package main

import (
	"fmt"
	"os"
)

func open(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return data, nil
}
