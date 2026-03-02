package main

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	binPath      = "./bin"
	binaryPrefix = "app-"
)

func discoverBinaries() []BinaryFile {
	dir, err := os.ReadDir(binPath)
	if err != nil {
		return nil
	}

	var binaries []BinaryFile
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		if !strings.HasPrefix(file.Name(), binaryPrefix) {
			continue
		}

		path := filepath.Join(binPath, file.Name())
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		// check executable flag
		if info.Mode()&0o111 == 0 {
			continue
		}

		binaries = append(binaries, BinaryFile{
			FullPath: path,
			Name:     file.Name(),
		})
	}

	return binaries
}

type BinaryFile struct {
	FullPath string
	Name     string
}
