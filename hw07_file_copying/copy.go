package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSameFileProvided      = errors.New("same file provided")
)

func progressBarLimit(info os.FileInfo, limit, offset int64) int {
	fileSize := info.Size()
	if (fileSize - offset) < limit {
		return int(fileSize - offset)
	}
	return int(limit)
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Check files before
	if fromPath == toPath {
		return ErrSameFileProvided
	}

	// Read source file
	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer sourceFile.Close()

	// Get file info
	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	// Set limit
	if limit <= 0 {
		limit = fileInfo.Size()
	}

	// Check offset limits
	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// Create target file
	targetFile, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer targetFile.Close()

	// Set offset
	_, err = sourceFile.Seek(offset, io.SeekStart)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	pbLimit := progressBarLimit(fileInfo, limit, offset)
	progressBar := pb.New(pbLimit)
	progressBar.SetUnits(pb.U_BYTES)
	progressBar.SetRefreshRate(time.Millisecond * 10)
	progressBar.Start()

	// Copy data
	reader := progressBar.NewProxyReader(sourceFile)
	_, err = io.CopyN(targetFile, reader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	progressBar.Finish()

	return nil
}
