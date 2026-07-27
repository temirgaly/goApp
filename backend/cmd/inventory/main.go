package main

import (
	"go.uber.org/fx"

	"github.com/orderpulse/backend/internal/inventory"
)

func main() {
	fx.New(
		inventory.Module,
	).Run()
}
