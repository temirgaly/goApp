package main

import (
	"go.uber.org/fx"

	"github.com/orderpulse/internal/worker"
)

func main() {
	fx.New(
		worker.Module,
	).Run()
}
