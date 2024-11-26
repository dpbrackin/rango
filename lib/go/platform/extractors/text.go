package extractors

import (
	"context"
	"io"
	"rango/core"
)

type TextExtractor struct {
	Storage core.StorageBackend
}

// Extract implements core.ContentExtractor.
func (t *TextExtractor) Extract(ctx context.Context, doc core.Document, w io.Writer) error {
	return t.Storage.Download(ctx, doc.Source, w)
}
