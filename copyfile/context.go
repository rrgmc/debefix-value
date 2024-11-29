package copyfile

import (
	"context"
	"errors"
)

// ToContext adds the Value context to the [context.Context].
func ToContext(ctx context.Context, process *Process) context.Context {
	return context.WithValue(ctx, key, process)
}

type keyType struct{}

var key keyType

func fromContext(ctx context.Context) (*Process, bool) {
	v, ok := ctx.Value(key).(*Process)
	return v, ok
}

func fromContextCheck(ctx context.Context) (*Process, error) {
	v, ok := fromContext(ctx)
	if !ok {
		return nil, errors.New("copyfile process was not initialized")
	}
	return v, nil
}
