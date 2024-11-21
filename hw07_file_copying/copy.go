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
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := validate(fromPath, toPath, offset)
	if err != nil {
		return err
	}

	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	srcInfo, err := src.Stat()
	if err != nil {
		return err
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

func validate(fromPath, toPath string, offset int64) error {
	src, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	dst, _ := os.Stat(toPath)

	if src.IsDir() || os.SameFile(src, dst) || !src.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > src.Size() {
		return ErrOffsetExceedsFileSize
	}

	return nil
}
