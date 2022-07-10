package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error while opening input file: %w", err)
	}
	defer fromFile.Close()

	limit, err = countRealLimit(fromPath, limit, offset)
	if err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error while creating result file: %w", err)
	}
	defer toFile.Close()

	fmt.Print("Copy ", limit, " bytes from ", fromPath, " to ", toPath, ": ")

	return progressCopyFromFile(fromFile, toFile, offset, limit)
}

func countRealLimit(fromFileName string, inputLimit int64, offset int64) (int64, error) {
	fileInfo, err := os.Stat(fromFileName)
	if err != nil {
		return 0, err
	}

	fileSize := fileInfo.Size()
	if fileSize <= 0 {
		return 0, ErrUnsupportedFile
	}

	copySize := fileSize - offset
	if copySize < 0 {
		return 0, ErrOffsetExceedsFileSize
	}

	if inputLimit <= 0 || inputLimit > copySize {
		return copySize, nil
	}

	return inputLimit, nil
}

func progressCopyFromFile(src *os.File, dst io.Writer, offset, limit int64) error {
	progress := newProgress(limit)

	if offset > 0 {
		_, err := src.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	bufferSize := int64(math.Ceil(float64(limit) / 10))
	maxBufferSize := int64(1 << 20)
	if bufferSize > maxBufferSize {
		bufferSize = maxBufferSize
	}

	for totalWritten := int64(0); totalWritten < limit; {
		written, err := io.CopyN(dst, src, bufferSize)
		if err != nil {
			if errors.Is(err, io.EOF) {
				_ = progress.finish()
				break
			}

			return err
		}

		totalWritten += written
		_ = progress.setCurrent(totalWritten)
	}

	return nil
}
