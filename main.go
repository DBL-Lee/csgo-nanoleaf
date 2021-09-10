package main

import (
	"fmt"
	"github.com/dank/go-csgsi"
)

func main(){
	game := csgsi.New(10)
	go func() {
		for state := range game.Channel {
			if state.Added != nil && state.Added.Round != nil {
				if state.Round.Bomb == "planted" {
					fmt.Println("planted")
				}
				if state.Round.Win_team != "" {
					if state.Round.Win_team == "T" {
						if state.Round.Bomb == "exploded" {
							fmt.Println("exploded T win")
						} else {
							fmt.Println("killed T win")
						}
					}else {
						if state.Round.Bomb == "defused" {
							fmt.Println("defused CT win")
						} else {
							fmt.Println("killed CT win")
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
								fmt.Printf("fired\n")
							} else if weapon.State == "reloading" && weapon.Ammo_clip == weapon.Ammo_clip_max {
								fmt.Printf("reloaded\n")
							}
						}
					}
				}
				// killed someone
				if player.State != nil && prevPlayer.State != nil {
					if player.State.Flashed != 0 {
						fmt.Println("flashed:", float64(player.State.Flashed)/ 255.0)
					}
					if player.State.Smoked != 0 {
						fmt.Println("smoke:", float64(player.State.Smoked)/ 255.0)
					}
					if player.State.Health < prevPlayer.State.Health {
						fmt.Println("healthchanged to", player.State.Health)
					}
					if player.State.Round_kills > prevPlayer.State.Round_kills {
						if player.State.Round_killhs > prevPlayer.State.Round_killhs {
							fmt.Println("headshot")
						} else {
							fmt.Println("kill")
						}
					}
				}
			}
		}
	}()
	game.Listen(":3000")
}
