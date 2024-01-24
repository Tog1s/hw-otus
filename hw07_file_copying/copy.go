package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// func checkFiles() {
// 	return
// }

// func checkOffset() {
// 	return
// }

func setLimit(limit int64, info os.FileInfo) int64 {
	if limit <= 0 {
		limit = info.Size()
	}
	return limit
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Read source file
	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer sourceFile.Close()

	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}
	limit = setLimit(limit, fileInfo)

	// Create target file
	targetFile, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer targetFile.Close()

	// Offset
	_, err = sourceFile.Seek(offset, io.SeekStart)
	if errors.Is(err, io.EOF) {
		err = nil
	}
	if err != nil {
		return err
	}

	// Copy data
	_, err = io.CopyN(targetFile, sourceFile, limit)
	if err != nil {
		return err
	}

	return nil
}
