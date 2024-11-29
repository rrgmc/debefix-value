package copyfile

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/rrgmc/debefix/v2"
)

// File source/destination: filename

type FilenameData struct {
	Filename   string
	IsAbsolute bool // whether the file name is relative (default) or absolute. If relative, it must be resolved by a FilenameProvider.
}

var _ FileSource = FilenameData{}
var _ FileDestination = FilenameData{}

// Filename is a path to a file.
func Filename(filename string, options ...FilenameOption) FilenameData {
	var optns filenameOptions
	for _, opt := range options {
		opt(&optns)
	}
	ret := FilenameData{
		Filename:   filename,
		IsAbsolute: optns.isAbsolute,
	}
	return ret
}

// AbsoluteFilename is an absolute path to a file.
func AbsoluteFilename(filename string, options ...FilenameOption) FilenameData {
	return Filename(filename, slices.Concat([]FilenameOption{WithFileNameIsAbsolute(true)}, options)...)
}

func (f FilenameData) ResolveSource(ctx context.Context, resolvedData *debefix.ResolvedData, fieldName string,
	values debefix.ValuesMutable, item Value) (FileReader, bool, error) {
	if f.IsAbsolute {
		return FileReaderFilename(f.Filename, f.Filename), true, nil
	}
	filename, err := getFilename(ctx, FileFieldSource, item.Info, f.Filename)
	if err != nil {
		return nil, false, fmt.Errorf("could not resolve filename: %w", err)
	}
	return FileReaderFilename(f.Filename, filename), true, nil
}

func (f FilenameData) ResolveDestination(ctx context.Context, resolvedData *debefix.ResolvedData, fieldName string,
	values debefix.ValuesMutable, item Value) (FileWriter, bool, error) {
	if f.IsAbsolute {
		return FileWriterFilename(f.Filename, f.Filename), true, nil
	}
	filename, err := getFilename(ctx, FileFieldDestination, item.Info, f.Filename)
	if err != nil {
		return nil, false, fmt.Errorf("could not resolve filename: %w", err)
	}
	return FileWriterFilename(f.Filename, filename), true, nil
}

type FilenameOption func(data *filenameOptions)

type filenameOptions struct {
	isAbsolute bool // whether the file name is relative (default) or absolute. If relative, it must be resolved by a FilenameProvider.
}

func WithFileNameIsAbsolute(isAbsolute bool) FilenameOption {
	return func(data *filenameOptions) {
		data.isAbsolute = isAbsolute
	}
}

// File source/destination: filename

type FilenameValueData struct {
	Filename   debefix.Value
	IsAbsolute bool // whether the file name is relative (default) or absolute. If relative, it must be resolved by a FilenameProvider.
}

var _ FileSource = FilenameValueData{}
var _ FileDestination = FilenameValueData{}

// FilenameValue is a path to a file.
func FilenameValue(filename debefix.Value, options ...FilenameOption) FilenameValueData {
	var optns filenameOptions
	for _, opt := range options {
		opt(&optns)
	}
	ret := FilenameValueData{
		Filename:   filename,
		IsAbsolute: optns.isAbsolute,
	}
	return ret
}

func (f FilenameValueData) ResolveSource(ctx context.Context, resolvedData *debefix.ResolvedData, fieldName string,
	values debefix.ValuesMutable, item Value) (FileReader, bool, error) {
	originalFilename, ok, err := f.getOriginalFilename(ctx, resolvedData, values)
	if err != nil {
		return nil, false, fmt.Errorf("could not resolve value: %w", err)
	}
	if !ok {
		return nil, false, nil
	}
	if f.IsAbsolute {
		return FileReaderFilename(originalFilename, originalFilename), true, nil
	}
	filename, err := getFilename(ctx, FileFieldSource, item.Info, originalFilename)
	if err != nil {
		return nil, false, fmt.Errorf("could not resolve filename: %w", err)
	}
	return FileReaderFilename(originalFilename, filename), true, nil
}

func (f FilenameValueData) ResolveDestination(ctx context.Context, resolvedData *debefix.ResolvedData, fieldName string,
	values debefix.ValuesMutable, item Value) (FileWriter, bool, error) {
	originalFilename, ok, err := f.getOriginalFilename(ctx, resolvedData, values)
	if err != nil {
		return nil, false, fmt.Errorf("could not resolve value: %w", err)
	}
	if !ok {
		return nil, false, nil
	}
	if f.IsAbsolute {
		return FileWriterFilename(originalFilename, originalFilename), true, nil
	}
	filename, err := getFilename(ctx, FileFieldDestination, item.Info, originalFilename)
	if err != nil {
		return nil, false, fmt.Errorf("could not resolve filename: %w", err)
	}
	return FileWriterFilename(originalFilename, filename), true, nil
}

func (f FilenameValueData) getOriginalFilename(ctx context.Context, resolvedData *debefix.ResolvedData,
	values debefix.ValuesMutable) (string, bool, error) {
	fv, ok, err := f.Filename.ResolveValue(ctx, resolvedData, values)
	if err != nil {
		return "", false, fmt.Errorf("could not resolve value: %w", err)
	}
	if !ok {
		return "", false, nil
	}
	originalFilename, ok := fv.(string)
	if !ok {
		return "", false, fmt.Errorf("resolved value is not a string (got %T)", fv)
	}

	return originalFilename, true, nil
}

// FilenameFormat formats a file name using [debefix.ValueFormat].
func FilenameFormat(format string, args ...any) FilenameValueData {
	return FilenameValue(debefix.ValueFormat(format, args...))
}

// FilenameFormat formats a file name using [debefix.ValueFormat].
func FilenameFormatOpt(format string, args []any, options ...FilenameOption) FilenameValueData {
	return FilenameValue(debefix.ValueFormat(format, args...), options...)
}

// FilenameFormatTemplate formats a file name using [debefix.ValueFormatTemplate].
func FilenameFormatTemplate(template string, args map[string]any, options ...FilenameOption) FilenameValueData {
	return FilenameValue(debefix.ValueFormatTemplate(template, args), options...)
}

// FileReader: filename

type FileReaderFilenameData struct {
	Filename string // the original file name (maybe a relative path)
	FilePath string // the full file path
}

var _ FileReader = FileReaderFilenameData{}
var _ FileFilename = FileReaderFilenameData{}

func FileReaderFilename(filename string, filePath string) FileReaderFilenameData {
	return FileReaderFilenameData{
		Filename: filename,
		FilePath: filePath,
	}
}

func (f FileReaderFilenameData) NewReader(ctx context.Context) (io.ReadCloser, error) {
	return f.openFile()
}

func (f FileReaderFilenameData) FileInfo(ctx context.Context) (os.FileInfo, bool, error) {
	file, err := f.openFile()
	if err != nil {
		return nil, false, err
	}
	defer file.Close()
	finfo, err := file.Stat()
	if err != nil {
		return nil, false, fmt.Errorf("error reading file info: %w", err)
	}
	return finfo, true, nil
}

func (f FileReaderFilenameData) GetFilename() string {
	return f.Filename
}

func (f FileReaderFilenameData) openFile() (*os.File, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file '%s': %w", f.FilePath, err)
	}
	return file, nil
}

// FileReader: bytes

type FileReaderBytesData struct {
	Data []byte
}

var _ FileReader = FileReaderBytesData{}

func FileReaderBytes(data []byte) FileReaderBytesData {
	return FileReaderBytesData{
		Data: data,
	}
}

func (f FileReaderBytesData) NewReader(ctx context.Context) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(f.Data)), nil
}

func (f FileReaderBytesData) FileInfo(ctx context.Context) (os.FileInfo, bool, error) {
	return nil, false, nil
}

// FileWriter: filename

type FileWriterFilenameData struct {
	Filename string // the original file name (maybe a relative path)
	FilePath string // the full file path
}

var _ FileWriter = FileWriterFilenameData{}
var _ FileFilename = FileWriterFilenameData{}

func FileWriterFilename(filename string, filePath string) FileWriterFilenameData {
	return FileWriterFilenameData{
		Filename: filename,
		FilePath: filePath,
	}
}

func (f FileWriterFilenameData) NewWriter(ctx context.Context) (io.WriteCloser, error) {
	return f.createFile()
}

func (f FileWriterFilenameData) GetFilename() string {
	return f.Filename
}

func (f FileWriterFilenameData) createFile() (*os.File, error) {
	file, err := os.Create(f.FilePath)
	if err != nil {
		return nil, fmt.Errorf("could not create file '%s': %w", f.FilePath, err)
	}
	return file, nil
}

// FileReader: cache

type FileReaderCacheData struct {
	fileSource          FileReader
	cached              *bytes.Buffer
	cachedFileInfo      os.FileInfo
	cachedError         error
	cachedFileInfoError error
}

var _ FileReader = (*FileReaderCacheData)(nil)

func FileReaderCache(fileSource FileReader) *FileReaderCacheData {
	return &FileReaderCacheData{
		fileSource: fileSource,
	}
}

func (f *FileReaderCacheData) NewReader(ctx context.Context) (io.ReadCloser, error) {
	// TODO implement me
	panic("implement me")
}

func (f *FileReaderCacheData) FileInfo(ctx context.Context) (os.FileInfo, bool, error) {
	// TODO implement me
	panic("implement me")
}
