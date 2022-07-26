package main

import (
	"errors"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOfCopyProcess         = errors.New("cant copy to a new file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {

	file, err := getFileToCopyWithOffset(fromPath, offset)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer file.Close()

	// offset больше, чем размер файла - невалидная ситуация
	fileStat, errStat := file.Stat()
	if errStat != nil {
		return ErrUnsupportedFile
	}
	fileSize := fileStat.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if errCopy := copyWithProcessBar(file, toPath, fileSize, limit); errCopy != nil {
		return ErrOfCopyProcess
	}

	return nil
}

func getFileToCopyWithOffset(fromPath string, offset int64) (*os.File, error) {
	file, errOpen := os.Open(fromPath)
	if errOpen != nil {
		return nil, errOpen
	}
	_, errSeek := file.Seek(offset, 0)
	if errSeek != nil {
		return nil, errSeek
	}

	return file, nil
}

func copyWithProcessBar(file *os.File, toPath string, fileSize int64, limit int64) error {
	barLimit := fileSize
	reader := io.LimitReader(file, barLimit)
	bar := pb.Full.Start64(barLimit)
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
	defer bar.Finish()

	newFile, errNew := os.Create(toPath)
	if errNew != nil {
		return errNew
	}
	defer newFile.Close()

	if limit == 0 {
		limit = fileSize
	}
	barReader := bar.NewProxyReader(reader)
	_, errCopy := io.CopyN(newFile, barReader, limit)
	if errCopy != nil && errCopy != io.EOF {
		return errCopy
	}

	return nil
}
