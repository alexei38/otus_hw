package main

import (
	"errors"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открываем файл только на чтение
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()
	stat, err := src.Stat()
	if err != nil {
		return err
	}

	if stat.Size() == 0 {
		return ErrUnsupportedFile
	}

	if stat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	// расчитываем limit с учетем offset до конца файла
	size := stat.Size() - offset
	if limit == 0 || limit > size {
		limit = size
	}

	// Сдвигаем указатель начала файла на offset
	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	// Создаем пустой dst файл и перезаписываем содержимое
	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Размер progress-bar - это количество байт, которые будем копировать
	bar := progressbar.DefaultBytes(
		limit,
		"copying",
	)

	// Пишем данные в dst файл и в progress bar с учетом лимита
	_, err = io.CopyN(io.MultiWriter(dst, bar), src, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
