package port

import (
	"context"

	"github.com/malikbenkirane/oha-opus-major/internal/domain/player"
)

// PlayerDataRepository abstracts the retrieval of player data.
// It allows different implementations (e.g., database, in‑memory, remote service)
// to be swapped without changing the underlying logic that depends on it.
type PlayerDataRepository interface {
	// Players returns all player records available in the data source.
	// The provided context should be used to honor cancellation and timeouts.
	// On success it returns a slice of player.Data; otherwise an error describing
	// the failure condition.
	Players(ctx context.Context) ([]player.Data, error)
}
