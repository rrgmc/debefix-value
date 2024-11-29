package copyfile

import (
	"context"
	"io"
	"os"

	"github.com/rrgmc/debefix/v2"
)

type FileSource interface {
	ResolveSource(ctx context.Context, resolvedData *debefix.ResolvedData, tableID debefix.TableID, fieldName string,
		values debefix.ValuesMutable, item Value) (FileReader, bool, error)
}

type FileDestination interface {
	ResolveDestination(ctx context.Context, resolvedData *debefix.ResolvedData, tableID debefix.TableID, fieldName string,
		values debefix.ValuesMutable, item Value) (FileWriter, bool, error)
}

type FileReader interface {
	NewReader(ctx context.Context) (io.ReadCloser, error)
	FileInfo(ctx context.Context) (os.FileInfo, bool, error)
	String() string
}

type FileWriter interface {
	NewWriter(ctx context.Context) (io.WriteCloser, error)
	String() string
}

type FileFilename interface {
	GetFilename() string
}
