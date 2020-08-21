package coffeemachine

import (
	"fmt"
)

// Machine maintains the configurations of a machine
type Machine struct {
	Outlet      int                     `json:"count_n"`
	Ingredients map[Ingredient]float64  `json:"ingredients"`
	Beverages   map[BeverageName]Recipe `json:"beverages"`
}

type Ingredient string

type BeverageName string

type Recipe map[Ingredient]float64

func (m *Machine) Menu() {
	fmt.Println("This machine can serve:")
	for key, _ := range m.Beverages {
		fmt.Printf("* %v\n", key)
	}
}

func (m *Machine) AddIngredient(ing Ingredient, quant float64) {
	m.Ingredients[ing] = quant
}

func (m *Machine) AddBeverage(name BeverageName, recipe Recipe) {
	m.Beverages[name] = recipe
}

func (m *Machine) Serving(bvg BeverageName) bool {
	_, found := m.Beverages[bvg]
	return found
}
