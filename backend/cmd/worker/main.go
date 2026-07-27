package main

import (
	"go.uber.org/fx"

	"github.com/orderpulse/backend/internal/worker"
)

func main() {
	fx.New(
		worker.Module,
	).Run()
}
