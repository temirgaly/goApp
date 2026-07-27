package main

import (
	"go.uber.org/fx"

	"github.com/orderpulse/backend/internal/order"
)

func main() {
	fx.New(
		order.Module,
	).Run()
}
