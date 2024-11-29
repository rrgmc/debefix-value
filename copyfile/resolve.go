package copyfile

import (
	"context"
	"fmt"
	"io"

	"github.com/rrgmc/debefix/v2"
)

// ResolveCallback is the callback to resolve the field values and/or copy the file.
type ResolveCallback func(ctx context.Context, resolvedData *debefix.ResolvedData, tableID debefix.TableID, fieldName string,
	values debefix.ValuesMutable, item Value, reader FileReader, writer FileWriter) error

// ResolveCopyFile is a ResolveCallback which copies the data from reader to writer.
func ResolveCopyFile(ctx context.Context, resolvedData *debefix.ResolvedData, tableID debefix.TableID, fieldName string,
	values debefix.ValuesMutable, item Value, reader FileReader, writer FileWriter) error {
	return CopyFile(ctx, reader, writer)
}

// CopyFile copies the data from reader to writer.
func CopyFile(ctx context.Context, reader FileReader, writer FileWriter) error {
	src, err := reader.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("error reading source: %w", err)
	}
	defer src.Close()

	dst, err := writer.NewWriter(ctx)
	if err != nil {
		return fmt.Errorf("error writing to destination: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("error copying data from source to destination: %w", err)
	}

	return nil
}
