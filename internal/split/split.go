package split

import (
	"errors"
	"fmt"
	"io"
	iofs "io/fs"
	"os"
	"path/filepath"

	"github.com/bodgit/gc"
	"github.com/bodgit/psx"
	"github.com/spf13/afero"
)

const maxChannels = 8

var fs = afero.NewOsFs()

var (
	errNoFreeChannels    = errors.New("no free memory card channels")
	errNotDirectory      = errors.New("not a directory")
	errUnknownMemoryCard = errors.New("unknown memory card")
)

func MemoryCards(dir string, files []string, useSize, useFlashID bool) error {
	base, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("unable to create absolute path: %w", err)
	}

	if fi, err := fs.Stat(base); err != nil || !fi.IsDir() {
		if err != nil {
			return fmt.Errorf("unable to stat directory: %w", err)
		}

		return errNotDirectory
	}

	for _, f := range files {
		source, err := filepath.Abs(f)
		if err != nil {
			return fmt.Errorf("unable to create absolute path: %w", err)
		}

		file, err := fs.Open(source)
		if err != nil {
			return fmt.Errorf("unable to open: %w", err)
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			return fmt.Errorf("unable to stat: %w", err)
		}

		if err := splitMemoryCard(base, file, fi.Size(), useSize, useFlashID); err != nil {
			return err
		}
	}

	return nil
}

type fileReader interface {
	io.Reader
	io.ReaderAt
}

func splitMemoryCard(base string, file fileReader, size int64, useSize, useFlashID bool) error {
	ok, err := psx.DetectMemoryCard(file, size)
	if err != nil {
		return fmt.Errorf("error detecting PlayStation memory card: %w", err)
	}

	if ok {
		return splitPSXMemoryCard(base, file)
	}

	ok, err = gc.DetectMemoryCard(file, size)
	if err != nil {
		return fmt.Errorf("error detecting GameCube memory card: %w", err)
	}

	if ok {
		return splitGCMemoryCard(base, file, useSize, useFlashID)
	}

	return errUnknownMemoryCard
}

func newMemoryCardFile(base, dir string) (io.WriteCloser, error) {
	directory := filepath.Join(base, dir)
dir:
	fi, err := fs.Stat(directory)

	if err != nil {
		if os.IsNotExist(err) {
			if err := fs.Mkdir(directory, os.ModePerm|os.ModeDir); err != nil {
				return nil, fmt.Errorf("unable to create directory: %w", err)
			}

			goto dir
		}

		return nil, fmt.Errorf("unable to stat directory: %w", err)
	}

	if !fi.IsDir() {
		return nil, errNotDirectory
	}

	var (
		i      int
		target string
	)

	for i = 1; i <= maxChannels; i++ {
		target = filepath.Join(directory, fmt.Sprintf("%s-%d.mcd", dir, i))
		if _, err = fs.Stat(target); err != nil {
			if os.IsNotExist(err) {
				break
			}

			return nil, fmt.Errorf("unable to stat file: %w", err)
		}
	}

	if i > maxChannels {
		return nil, errNoFreeChannels
	}

	file, err := fs.Create(target)
	if err != nil {
		return nil, fmt.Errorf("unable to create file: %w", err)
	}

	return file, nil
}

type cardReader interface {
	Open() (iofs.File, error)
}

type cardWriter interface {
	Create() (io.WriteCloser, error)
}

func copyData(r cardReader, w cardWriter) error {
	rc, err := r.Open()
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer rc.Close()

	wc, err := w.Create()
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}

	if _, err := io.Copy(wc, rc); err != nil {
		return fmt.Errorf("unable to copy: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("unable to close: %w", err)
	}

	return nil
}
