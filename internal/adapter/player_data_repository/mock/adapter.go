package mock

import (
	"context"

	"github.com/malikbenkirane/oha-opus-major/internal/domain/player"
	"github.com/malikbenkirane/oha-opus-major/internal/domain/weapon"
	"github.com/malikbenkirane/oha-opus-major/internal/port"
)

func New() port.PlayerDataRepository {
	return &adapter{}
}

type adapter struct{}

func (a adapter) Players(ctx context.Context) ([]player.Data, error) {
	return []player.Data{
		{
			Position: player.Vector2{},
			Weapons: []weapon.Weapon{
				{
					Name: "butterfly",
				},
			},
		},
	}, nil
}
