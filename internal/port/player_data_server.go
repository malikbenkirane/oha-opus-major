package port

import "context"

type PlayerDataServer interface {
	Serve(ctx context.Context) error
}
