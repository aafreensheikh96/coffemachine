package coffeemachine

import (
	"context"
)

type Service interface {
	// Menu displays the menu of the machine
	Menu()

	// TakeOrders, take orders to be served
	TakeOrders(ctx context.Context, bvgName BeverageName, orderNo int64)

	// IngredientsAvailiblity checks how much ingreients are still available based on the recipes.
	IngredientsAvailiblity(ctx context.Context)

	// RefillIngredient can refill ingredients.
	RefillIngredient(ctx context.Context, ing Ingredient, quant float64)

	// AddRecipe, can add a new recipe to the machine.
	AddRecipe(ctx context.Context, bvgName BeverageName, recipe Recipe)

	// checks if all orders are fulfilled and closes
	AllDone(ctx context.Context)
}
