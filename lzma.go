package archives

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/ulikunitz/xz/lzma"
)

func init() {
	RegisterFormat(Lzma{})
}

// Lzma facilitates lzma compression.
type Lzma struct {
}

func (Lzma) Extension() string { return ".lzma" }
func (Lzma) MediaType() string { return "application/x-lzma" }

func (lzma Lzma) Match(_ context.Context, filename string, stream io.Reader) (MatchResult, error) {
	var mr MatchResult

	// match filename
	if strings.Contains(strings.ToLower(filename), lzma.Extension()) {
		mr.ByName = true
	}

	// match file header
	buf, err := readAtMost(stream, len(lzmaHeader))
	if err != nil {
		return mr, err
	}
	mr.ByStream = bytes.Equal(buf, lzmaHeader)

	return mr, nil
}

func (Lzma) OpenWriter(r io.Writer) (io.WriteCloser, error) {
	lzmaWriter, err := lzma.NewWriter(r)
	return lzmaWriter, err
}

func (Lzma) OpenReader(r io.Reader) (io.ReadCloser, error) {
	lzmaReader, err := lzma.NewReader(r)
	return io.NopCloser(lzmaReader), err
}

// magic number at the beginning of lzma files
var lzmaHeader = []byte{0x5d, 0x00, 0x00}
