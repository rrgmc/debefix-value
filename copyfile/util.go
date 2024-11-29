package copyfile

import (
	"context"
	"errors"

	"github.com/rrgmc/debefix/v2"
)

func getFilename(ctx context.Context, fileField FileField, item Value, tableID debefix.TableID, filename string) (string, error) {
	process, err := fromContextCheck(ctx)
	if err != nil {
		return "", err
	}
	if process.filenameProvider == nil {
		return "", errors.New("filename provider not set")
	}

	retfilename, err := process.filenameProvider(ctx, fileField, item, tableID, filename)
	if err != nil {
		return "", err
	}

	return retfilename, nil
}
