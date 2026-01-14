# A 2‑D Multiplayer Game Service Model

Only the 

    GET /update-data-player

endpoint is needed; it provides the current state of every player in the game.

## Game Overview

Let's say we've built a 2‑D demo where every player has a position, an
orientation, and the ability to fire weapons.

    type Vector2 struct {
        X, Y float32
    }

    // Player represents a 2‑D game entity with a position, orientation, and
    // the ability to fire weapons.
    type Player struct {
        Position    Vector2 // (x, y) coordinates in world space
        Orientation float64 // angle in radians, where 0 points to the right
        CanFire     bool    // true if the player is currently able to fire
        Weapons     []Weapon // Weapons holds the character's equipped weapons.
	    Health      float64 // Health is a percentage between 0 and 1.
    }

    type Weapon struct {
        Name   string  // weapon identifier
        Damage float64 // fixed damage applied on hit
    }

### Weapons Definition

Each weapon is identified by a name and a fixed damage value applied when it
hits the targeted player.

    type Weapon struct {
        Name   string  // weapon identifier
        Damage float64 // fixed damage applied on hit
    }

### Potential Velocity Enhancement

Adding a projectile speed would make gameplay richer: players could dodge
attacks that haven’t yet reached them, introducing timing and positioning
skills.  

Projectile physics, collision prediction, and latency handling are excluded
from this homework assignment.

### Weapon Switching

Players can carry several weapons and switch among them during play.

    // Weapons holds the character's equipped weapons.
    Weapons     []Weapon


### Health System

Players possess health points and can restore them by collecting health items placed in the level.  

	// Health is a percentage between 0 and 1.
	Health      float64


### World Modeling
The environment is not formally modeled; this avoids added complexity that is out of scope.


## Trade-offs

Since the project’s primary emphasis is on DevOps rather than full‑scale game
development, we’ve deliberately narrowed the scope to basic player data.

- **Gameplay scope:** Limited to movement and firing only
- **Weapon data:** Single `damage` value per weapon
- **Projectile physics:** No velocity or advanced physics
- **World model:** No tiles, obstacles, or terrain:
  entirely centered on player data
