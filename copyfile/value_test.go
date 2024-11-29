package copyfile

import (
	"context"
	"testing"

	"github.com/rrgmc/debefix/v2"
	"gotest.tools/v3/assert"
)

func TestValue(t *testing.T) {
	process := NewProcess(
		WithProcessResolveCallback(func(ctx context.Context, resolvedData *debefix.ResolvedData, tableID debefix.TableID,
			fieldName string, values debefix.ValuesMutable, item Value, reader FileReader, writer FileWriter) error {
			// fi, ok, err := reader.FileInfo(ctx)
			// if err != nil {
			// 	return err
			// }
			// if ok {
			// 	values.Set("file_size", fi.Size())
			// }
			return nil
		}),
	)
	ctx, err := process.Start(context.Background())
	assert.NilError(t, err)

	rd := debefix.NewResolvedData()

	// v := New(nil, Filename("x/y.yaml"), Filename("h/b.yaml"))
	v := New(nil,
		AbsoluteFilename("x/y.yaml"),
		AbsoluteFilename("h/b.yaml"))

	values := debefix.MapValues{}

	err = v.Resolve(ctx, rd, debefix.TableName("x"), "_copyfile", values)
	assert.NilError(t, err)

	// fmt.Println(v)
}
