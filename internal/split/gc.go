package split

import (
	"fmt"
	"io"

	"github.com/bodgit/gc"
)

const extensionGC = "raw"

// Some games report a different game + publisher code on the memory card
// versus the game disc.
var mcCodeToDiscCode = map[string]string{
	"GFZE8P": "GFZE01", // F-Zero GX
	"GFZJ8P": "GFZJ01",
	"GFZP8P": "GFZP01",
}

func sanitizeGCCode(gameCode, makerCode string) string {
	code := gameCode + makerCode

	if newCode, ok := mcCodeToDiscCode[code]; ok {
		return newCode
	}

	return code
}

//nolint:cyclop
func splitGCMemoryCard(base string, fr io.Reader, useSize, useFlashID bool) error {
	r, err := gc.NewReader(fr)
	if err != nil {
		return fmt.Errorf("unable to create GameCube reader: %w", err)
	}

	codes := make(map[string]struct{})
	for _, file := range r.File {
		codes[sanitizeGCCode(file.GameCode, file.MakerCode)] = struct{}{}
	}

	for code := range codes {
		f, err := newMemoryCardFile(base, code, extensionGC)
		if err != nil {
			return err
		}
		defer f.Close()

		cardSize := gc.MemoryCard251
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
			if sanitizeGCCode(file.GameCode, file.MakerCode) != code {
				continue
			}

			if err := copyData(file, w); err != nil {
				return err
			}
		}
	}

	return nil
}
