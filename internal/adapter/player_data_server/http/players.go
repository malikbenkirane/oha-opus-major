package http

import "github.com/malikbenkirane/oha-opus-major/internal/domain/player"

// Builds a JSON payload for player data.
//
// Ensure that *playerDataJSON instance is allocated beforehand
// to prevent a nil pointer dereference.
func (p *playerDataJSON) from(playerData player.Data) {
	weapons := make([]weaponDataJSON, len(playerData.Weapons))
	for i, weapon := range playerData.Weapons {
		weapons[i] = weaponDataJSON{
			Name:   weapon.Name,
			Damage: weapon.Damage,
		}
	}
	p.Weapons = weapons
	p.PositionX = playerData.Position.X
	p.PositionY = playerData.Position.Y
	p.Orientation = playerData.Orientation
	p.CanFire = playerData.CanFire
	p.Health = playerData.Health
}

type playerDataJSON struct {
	PositionX   float64
	PositionY   float64
	Orientation float64
	CanFire     bool
	Weapons     []weaponDataJSON
	Health      float64
}

type weaponDataJSON struct {
	Name   string
	Damage float64
}
