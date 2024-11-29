package copyfile

import (
	"context"
	"errors"
)

func getFilename(ctx context.Context, fileField FileField, info any, filename string) (string, error) {
	process, err := fromContextCheck(ctx)
	if err != nil {
		return "", err
	}
	if process.filenameProvider == nil {
		return "", errors.New("filename provider not set")
	}

	retfilename, err := process.filenameProvider(ctx, fileField, info, filename)
	if err != nil {
		return "", err
	}

	return retfilename, nil
}
