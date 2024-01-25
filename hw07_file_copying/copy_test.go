package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Input equal output", func(t *testing.T) {
		inputFile := "testdata/input.txt"
		outputFile := os.TempDir() + "/out.txt"
		err := Copy(inputFile, outputFile, 0, 0)

		sourceFile, _ := os.Open(inputFile)
		sourceFileInfo, _ := sourceFile.Stat()

		targetFile, _ := os.Open(outputFile)
		targetFileInfo, _ := targetFile.Stat()

		require.Equal(t, nil, err)
		require.Equal(t, sourceFileInfo.Size(), targetFileInfo.Size())
	})

	t.Run("same file return error", func(t *testing.T) {
		inputFile := "testdata/input.txt"

		err := Copy(inputFile, inputFile, 0, 0)
		require.Equal(t, ErrSameFileProvided, err)
	})

	t.Run("nonexistent file return error", func(t *testing.T) {
		testFile := os.TempDir() + "/nonexistent_file"
		err := Copy(testFile, os.TempDir(), 0, 1024)

		require.Equal(t, ErrUnsupportedFile, err)
	})
}
