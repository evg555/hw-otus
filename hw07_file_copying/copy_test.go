package main

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fileDir           = "testdata/"
	inputFile         = "input.txt"
	outputFile        = "out.txt"
	notDefinedLenFile = "/dev/urandom"
	notExistFile      = "not_exist"
)

func TestCopy(t *testing.T) {
	t.Run("full copying", func(t *testing.T) {
		out, err := os.CreateTemp(fileDir, outputFile)
		if err != nil {
			t.Error(err)
			return
		}
		defer os.Remove(out.Name())

		err = Copy(fileDir+inputFile, out.Name(), 0, 0)
		require.Nil(t, err)
		require.FileExists(t, out.Name())

		inputInfo, _ := os.Stat(fileDir + inputFile)
		outInfo, _ := out.Stat()

		require.Equal(t, inputInfo.Size(), outInfo.Size())
	})

	t.Run("copying with limit less than len of file", func(t *testing.T) {
		limit = 1000

		out, err := os.CreateTemp(fileDir, outputFile)
		if err != nil {
			t.Error(err)
			return
		}
		defer os.Remove(out.Name())

		err = Copy(fileDir+inputFile, out.Name(), 0, limit)
		require.Nil(t, err)
		require.FileExists(t, out.Name())

		outInfo, _ := out.Stat()

		require.LessOrEqual(t, limit, outInfo.Size())
	})

	t.Run("copying with limit more than len of file", func(t *testing.T) {
		limit = 10000

		out, err := os.CreateTemp(fileDir, outputFile)
		if err != nil {
			t.Error(err)
			return
		}
		defer os.Remove(out.Name())

		err = Copy(fileDir+inputFile, out.Name(), 0, limit)
		require.Nil(t, err)
		require.FileExists(t, out.Name())

		inputInfo, _ := os.Stat(fileDir + inputFile)
		outInfo, _ := out.Stat()

		require.Equal(t, inputInfo.Size(), outInfo.Size())
	})

	t.Run("copying with offset", func(t *testing.T) {
		offset = 1000

		out, err := os.CreateTemp(fileDir, outputFile)
		if err != nil {
			t.Error(err)
			return
		}
		defer os.Remove(out.Name())

		err = Copy(fileDir+inputFile, out.Name(), offset, 0)
		require.Nil(t, err)
		require.FileExists(t, out.Name())

		inputInfo, _ := os.Stat(fileDir + inputFile)
		outInfo, _ := out.Stat()

		require.Equal(t, inputInfo.Size()-offset, outInfo.Size())
	})

	t.Run("copying with offset and limit", func(t *testing.T) {
		offset = 100
		limit = 1000

		out, err := os.CreateTemp(fileDir, outputFile)
		if err != nil {
			t.Error(err)
			return
		}
		defer os.Remove(out.Name())

		err = Copy(fileDir+inputFile, out.Name(), offset, limit)
		require.Nil(t, err)
		require.FileExists(t, out.Name())

		outInfo, _ := out.Stat()

		require.Equal(t, limit, outInfo.Size())
	})

	t.Run("copying with offset and limit more than len of file", func(t *testing.T) {
		offset = 6000
		limit = 1000

		out, err := os.CreateTemp(fileDir, outputFile)
		if err != nil {
			t.Error(err)
			return
		}
		defer os.Remove(out.Name())

		err = Copy(fileDir+inputFile, out.Name(), offset, limit)
		require.Nil(t, err)
		require.FileExists(t, out.Name())

		inInfo, _ := os.Stat(fileDir + inputFile)
		outInfo, _ := out.Stat()

		fmt.Println(inInfo.Size() - offset)

		require.Equal(t, inInfo.Size()-offset, outInfo.Size())
	})

	t.Run("overwriting existed out file", func(t *testing.T) {
		limit = 10

		out, err := os.CreateTemp(fileDir, outputFile)
		if err != nil {
			t.Error(err)
			return
		}
		defer os.Remove(out.Name())

		err = Copy(fileDir+inputFile, out.Name(), 0, 0)
		require.Nil(t, err)
		require.FileExists(t, out.Name())

		inputInfo, _ := os.Stat(fileDir + inputFile)
		outInfo, _ := out.Stat()

		require.Equal(t, inputInfo.Size(), outInfo.Size())

		err = Copy(fileDir+inputFile, out.Name(), 0, limit)
		require.Nil(t, err)
		require.FileExists(t, out.Name())

		outInfo, _ = out.Stat()

		require.Equal(t, limit, outInfo.Size())
	})
}

func TestCopyError(t *testing.T) {
	t.Run("got directory instead file", func(t *testing.T) {
		err := Copy(fileDir, "", 0, 0)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrIsDirectory)
	})

	t.Run("file not exist", func(t *testing.T) {
		err := Copy(fileDir+notExistFile, "", 0, 0)
		require.Error(t, err)
		require.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("unsupported file", func(t *testing.T) {
		err := Copy(notDefinedLenFile, "", 0, 0)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("out file not defined", func(t *testing.T) {
		err := Copy(fileDir+inputFile, "", 0, 0)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("offset more than len of file", func(t *testing.T) {
		offset = 7000

		err := Copy(fileDir+inputFile, outputFile, offset, 0)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
}
