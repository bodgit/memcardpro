package split

import (
	"fmt"
	"io"

	"github.com/bodgit/gc"
)

func splitGCMemoryCard(base string, fr io.Reader, retainSize bool) error {
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
		if retainSize {
			cardSize = r.CardSize
		}

		w, err := gc.NewWriter(f, cardSize, r.Encoding)
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
