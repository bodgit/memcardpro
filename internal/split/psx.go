package split

import (
	"fmt"
	"io"

	"github.com/bodgit/psx"
)

const extensionPSX = "mcd"

func sanitizePSXCode(code string) string {
	if code[4] == 'P' {
		return code[:4] + "-" + code[5:]
	}

	return code
}

func splitPSXMemoryCard(base string, fr io.Reader) error {
	r, err := psx.NewReader(fr)
	if err != nil {
		return fmt.Errorf("unable to create PlayStation reader: %w", err)
	}

	codes := make(map[string]struct{})
	for _, file := range r.File {
		codes[sanitizePSXCode(file.ProductCode)] = struct{}{}
	}

	for code := range codes {
		f, err := newMemoryCardFile(base, code, extensionPSX)
		if err != nil {
			return err
		}
		defer f.Close()

		w, err := psx.NewWriter(f)
		if err != nil {
			return fmt.Errorf("unable to create PlayStation writer: %w", err)
		}
		defer w.Close()

		for _, file := range r.File {
			if sanitizePSXCode(file.ProductCode) != code {
				continue
			}

			if err := copyData(file, w); err != nil {
				return err
			}
		}
	}

	return nil
}
