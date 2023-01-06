package split

import (
	"fmt"
	"io"

	"github.com/bodgit/gc"
)

//nolint:cyclop
func splitGCMemoryCard(base string, fr io.Reader, useSize, useFlashID bool) error {
	r, err := gc.NewReader(fr)
	if err != nil {
		return fmt.Errorf("unable to create GameCube reader: %w", err)
	}

	codes := make(map[string]struct{})
	for _, file := range r.File {
		codes[file.GameCode] = struct{}{}
	}

	for code := range codes {
		f, err := newMemoryCardFile(base, code)
		if err != nil {
			return err
		}
		defer f.Close()

		cardSize := gc.MemoryCard2043
		if useSize {
			cardSize = r.CardSize
		}

		options := []func(*gc.Writer) error{
			gc.CardSize(cardSize),
			gc.Encoding(r.Encoding),
		}

		if useFlashID {
			options = append(options, gc.FlashID(r.FlashID))
		} else {
			options = append(options, gc.FormatTime(0))
		}

		w, err := gc.NewWriter(f, options...)
		if err != nil {
			return fmt.Errorf("unable to create GameCube writer: %w", err)
		}
		defer w.Close()

		for _, file := range r.File {
			if file.GameCode != code {
				continue
			}

			if err := copyData(file, w); err != nil {
				return err
			}
		}
	}

	return nil
}
