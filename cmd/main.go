package main

import (
	"context"

	"github.com/DBL-Lee/csgo-nanoleaf/internal/csgo"
	"github.com/DBL-Lee/csgo-nanoleaf/internal/nanoleaf"
	"github.com/dank/go-csgsi"
)

func main() {
	game := csgsi.New(10)
	events := csgo.NewCSGOServer(game)
	n := nanoleaf.NewNanoLeaf(
		events,
		"http://192.168.1.217:16021",
		[]string{},
		"tTbNozk7tezDTUjtPtWIOLrHXoo6TJ9I",
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go n.Start(ctx)
	game.Listen(":3000")
}
