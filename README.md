# coffemachine

Implementation of a coffemachine

# Setup
`go run go run ./cmd/coffeemachine/main.go`

pass a mchine and an order file name as params.
`handler.InitialSetup("machine1", "orders1")` in the main

This is a working coffee machine example.

The machine has functionailities:
  Machine is configurable, based on number of outlets, ingredients it can contain and beverages it can served based on respective recipe.
  Can add a recipe
  Can add/refill in any ingredient
  Can take multiple orders
  Indicate if a ingredient is running low
  