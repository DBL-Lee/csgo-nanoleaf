package csgo

import (
	"fmt"

	"github.com/dank/go-csgsi"
)

type Events struct {
	BombPlanted  chan struct{}
	BombExploded chan struct{}
	BombDefused  chan struct{}
	RoundEnd     chan struct{}
	Fire         chan struct{}
	Kill         chan struct{}
	Headshot     chan struct{}
	Flashed      chan float64
	Smoked       chan float64
	Health       chan int
}

func NewCSGOServer(game *csgsi.Game) *Events {
	e := &Events{
		BombPlanted:  make(chan struct{}),
		BombExploded: make(chan struct{}),
		BombDefused:  make(chan struct{}),
		RoundEnd:     make(chan struct{}),
		Fire:         make(chan struct{}),
		Kill:         make(chan struct{}),
		Headshot:     make(chan struct{}),
		Flashed:      make(chan float64, 100),
		Smoked:       make(chan float64, 100),
		Health:       make(chan int, 100),
	}
	go e.Start(game)
	return e
}

func (e *Events) Start(game *csgsi.Game) {
	for state := range game.Channel {
		if state.Added != nil && state.Added.Round != nil {
			if state.Round.Bomb == "planted" {
				select {
				case e.BombPlanted <- struct{}{}:
				default:
				}
			}
			if state.Round.Win_team != "" {
				if state.Round.Win_team == "T" {
					if state.Round.Bomb == "exploded" {
						select {
						case e.BombExploded <- struct{}{}:
						default:
						}
					} else {
						select {
						case e.RoundEnd <- struct{}{}:
						default:
						}
					}
				} else {
					if state.Round.Bomb == "defused" {
						select {
						case e.BombDefused <- struct{}{}:
						default:
						}
					} else {
						select {
						case e.RoundEnd <- struct{}{}:
						default:
						}
					}
				}
			}
		}
		if state.Player != nil && state.Previously != nil && state.Previously.Player != nil {
			player := state.Player
			prevPlayer := state.Previously.Player
			for k, weapon := range state.Player.Weapons {
				prevWeapon := state.Previously.Player.Weapons[k]
				if prevWeapon != nil {
					// prevWeapon name is not empty only when
					// changing weapon
					if prevWeapon.Name == "" {
						if weapon.Ammo_clip < prevWeapon.Ammo_clip {
							select {
							case e.Fire <- struct{}{}:
							default:
							}
						} else if weapon.State == "reloading" && weapon.Ammo_clip == weapon.Ammo_clip_max {
							fmt.Printf("reloaded\n")
						}
					}
				}
			}
			// killed someone
			if player.State != nil && prevPlayer.State != nil {
				if player.State.Flashed != 0 {
					select {
					case e.Flashed <- float64(player.State.Flashed) / 255.0:
					default:
					}
				}
				if player.State.Smoked != 0 {
					select {
					case e.Smoked <- float64(player.State.Smoked) / 255.0:
					default:
					}
				}
				if player.State.Health < prevPlayer.State.Health {
					select {
					case e.Health <- player.State.Health:
					default:
					}
				}
				if player.State.Round_kills > prevPlayer.State.Round_kills {
					if player.State.Round_killhs > prevPlayer.State.Round_killhs {
						select {
						case e.Headshot <- struct{}{}:
						default:
						}
					} else {
						select {
						case e.Kill <- struct{}{}:
						default:
						}
					}
				}
			}
		}
	}
}
