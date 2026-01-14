package player

import "github.com/malikbenkirane/oha-opus-major/internal/domain/weapon"

// Player represents a 2‑D game entity with a position, orientation, and
// the ability to fire weapons.
type Data struct {
	Position    Vector2         // (x, y) coordinates in world space
	Orientation float64         // angle in radians, where 0 points to the right
	CanFire     bool            // true if the player is currently able to fire
	Weapons     []weapon.Weapon // Weapons holds the character's equipped weapons.
	Health      float64         // Health is a percentage between 0 and 1.
}

type Vector2 struct {
	X, Y float64
}
