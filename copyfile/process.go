package copyfile

import (
	"context"

	"github.com/rrgmc/debefix/v2"
)

type FileField int

const (
	FileFieldSource FileField = iota
	FileFieldDestination
)

type FilenameProvider func(ctx context.Context, fileField FileField, item Value, tableID debefix.TableID, filename string) (string, error)

// Process is a [debefix.Process] that stores files to be copied, and copy them at the end of the process.
type Process struct {
	filenameProvider FilenameProvider
	resolveCallback  ResolveCallback // callback that can be used to set fields and/or copy the file
}

var _ debefix.Process = (*Process)(nil)

func NewProcess(options ...ProcessOption) debefix.Process {
	ret := &Process{}
	for _, option := range options {
		option(ret)
	}
	return ret
}

func (p *Process) Start(ctx context.Context) (context.Context, error) {
	ctx = ToContext(ctx, p)
	return ctx, nil
}

func (p *Process) Finish(ctx context.Context) error {
	return nil
}

type ProcessOption func(*Process)

func WithProcessResolveCallback(resolveCallback ResolveCallback) ProcessOption {
	return func(p *Process) {
		p.resolveCallback = resolveCallback
	}
}

func WithProcessFilenameProvider(filenameProvider FilenameProvider) ProcessOption {
	return func(p *Process) {
		p.filenameProvider = filenameProvider
	}
}
