package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrIsDirectory           = errors.New("it is a directory")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer src.Close()

	srcInfo, err := src.Stat()
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return ErrIsDirectory
	}

	if srcInfo.Size() == 0 || toPath == "" {
		return ErrUnsupportedFile
	}

	if offset > srcInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if offset > 0 {
		_, err = src.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	if limit == 0 || limit+offset > srcInfo.Size() {
		limit = srcInfo.Size() - offset
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(src)
	defer barReader.Close()

	_, err = io.CopyN(dst, barReader, limit)
	if err != nil {
		return err
	}

	return nil
}
