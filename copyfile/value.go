package copyfile

import (
	"context"
	"errors"
	"fmt"

	"github.com/rrgmc/debefix/v2"
)

// Value is a [debefix.ValueMultiple] that can be used to copy files and optionally sets row values based on
// the file.
// Files are not handled by it, the copying itself must be done in a callback set by WithResolveCallback.
type Value struct {
	Info            any             // custom data
	Source          FileSource      // Source file
	Destination     FileDestination // Destination file
	resolveCallback ResolveCallback // callback that can be used to copy the file
}

var _ debefix.ValueMultiple = Value{}
var _ debefix.ValueDependencies = Value{}

// New creates a value for a field representing a file to be copied.
func New(info any, source FileSource, destination FileDestination, options ...ValueOption) Value {
	ret := Value{
		Info:        info,
		Source:      source,
		Destination: destination,
	}
	for _, option := range options {
		option(&ret)
	}
	return ret
}

func (v Value) Resolve(ctx context.Context, resolvedData *debefix.ResolvedData, tableID debefix.TableID, fieldName string, values debefix.ValuesMutable) error {
	reader, readerOk, err := v.Source.ResolveSource(ctx, resolvedData, tableID, fieldName, values, v)
	if err != nil {
		return fmt.Errorf("error resolving source value: %w", err)
	}
	if !readerOk {
		return debefix.ResolveLater
	}

	writer, writerOk, err := v.Destination.ResolveDestination(ctx, resolvedData, tableID, fieldName, values, v)
	if err != nil {
		return fmt.Errorf("error resolving destination value: %w", err)
	}
	if !writerOk {
		return debefix.ResolveLater
	}

	if v.resolveCallback != nil {
		err = v.resolveCallback(ctx, resolvedData, tableID, fieldName, values, v, reader, writer)
		if err != nil {
			return err
		}
	} else {
		process, ok := fromContext(ctx)
		if !ok || process.resolveCallback == nil {
			return errors.New("no callback found to process copy file result")
		}
		err = process.resolveCallback(ctx, resolvedData, tableID, fieldName, values, v, reader, writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v Value) TableDependencies() []debefix.TableID {
	var deps []debefix.TableID
	for _, vv := range []any{v.Source, v.Destination} {
		if vd, ok := vv.(debefix.ValueDependencies); ok {
			deps = append(deps, vd.TableDependencies()...)
		}
	}
	return deps
}

type ValueOption func(*Value)

// WithResolveCallback sets a callback to be used to effectively copy the file.
func WithResolveCallback(resolveCallback ResolveCallback) ValueOption {
	return func(v *Value) {
		v.resolveCallback = resolveCallback
	}
}
