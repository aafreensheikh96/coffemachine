package service

import (
	"coffeemachine"
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

// impl implements all componets for a machine service
type impl struct {
	machine          coffeemachine.Machine
	orderChan        chan coffeemachine.BeverageName
	wg               sync.WaitGroup
	ingredientAccess sync.Mutex
	outputFile       *os.File
}

// New starts a new machine, with the provided config
func New(m coffeemachine.Machine, outFile string) coffeemachine.Service {
	orderChan := make(chan coffeemachine.BeverageName, m.Outlet)
	var (
		ingLock sync.Mutex
		wg      sync.WaitGroup
	)
	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	return &impl{
		machine:          m,
		orderChan:        orderChan,
		ingredientAccess: ingLock,
		wg:               wg,
		outputFile:       f,
	}
}

func (svc *impl) TakeOrders(ctx context.Context, bvgName coffeemachine.BeverageName, orderNo int64) {
	// takes the incoming order and places into the `orderChan` channel
	select {
	case <-ctx.Done():
		fmt.Println("context done")
		return
	case svc.orderChan <- bvgName:
	}
	svc.wg.Add(1)
	go svc.orderWorker(ctx, orderNo)
}

func (svc *impl) AllDone(ctx context.Context) {
	svc.wg.Wait()
	fmt.Println("done!")
	close(svc.orderChan)
}

// orderWorker, it is worker to prepare orders
func (svc *impl) orderWorker(ctx context.Context, orderNo int64) {
	defer svc.wg.Done()

	// accepts the incoming order from the `orderChan` and quits if the context of the program is done.
	select {
	case <-ctx.Done():
		// fmt.Println("context done")
		return
	case bvg, ok := <-svc.orderChan:
		if !ok {
			// fmt.Println("not ok")
			return
		}
		if !svc.machine.Serving(bvg) {
			svc.print(fmt.Sprintf("Order: %v, %v is not served on this machine\n", orderNo, bvg))
			return
		}
		// `ingredientAccess` at any point only one go routine can have access to the ingredients
		// to avoid any race conditions
		svc.ingredientAccess.Lock()
		svc.prepareBeverage(ctx, bvg, orderNo)
		svc.ingredientAccess.Unlock()
		// svc.IngredientsAvailiblity(ctx)
	}
}

// prepareBeverage, prepares a `beverage` based on the inbuilt `recipe`
func (svc *impl) prepareBeverage(ctx context.Context, bvg coffeemachine.BeverageName, orderNo int64) {
	for k, v := range svc.machine.Beverages[bvg] {
		ingQuant, found := svc.machine.Ingredients[k]
		if !found {
			svc.print(fmt.Sprintf("Order: %v, %v cannot be prepared because item %v is not available\n", orderNo, bvg, k))
			return
		}
		if ingQuant < v {
			svc.print(fmt.Sprintf("Order: %v, %v cannot be prepared because item %v is not sufficient\n", orderNo, bvg, k))
			return
		}
	}
	for k, v := range svc.machine.Beverages[bvg] {
		ingQuant, _ := svc.machine.Ingredients[k]
		svc.machine.Ingredients[k] = ingQuant - v
	}

	svc.print(fmt.Sprintf("Order: %v, %v is prepared at %v\n", orderNo, bvg, time.Now().Format(time.Kitchen)))

}

func (svc *impl) IngredientsAvailiblity(ctx context.Context) {
	for ing, avlb := range svc.machine.Ingredients {
		// requiredQuant sum of a particular ingredient in all beverages,
		// if it is greater than the actual ingredient value present,
		// we indicate that ingredient will require a refill.
		requiredQuant := 0.0
		for _, rcp := range svc.machine.Beverages {
			quant, found := rcp[ing]
			if found {
				requiredQuant += quant
			}
		}
		if requiredQuant > avlb {
			svc.print(fmt.Sprintf("%v is running low, would require refill\n", ing))
		}
	}

}

func (svc *impl) RefillIngredient(ctx context.Context, ing coffeemachine.Ingredient, quant float64) {
	_, found := svc.machine.Ingredients[ing]
	if !found {
		svc.machine.Ingredients[ing] = quant
	} else {
		svc.machine.Ingredients[ing] += quant
	}
}

func (svc *impl) AddRecipe(ctx context.Context, bvgName coffeemachine.BeverageName, recipe coffeemachine.Recipe) {
	svc.machine.Beverages[bvgName] = recipe
}

func (svc *impl) Menu() {
	svc.machine.Menu()
}

func (svc *impl) print(str string) {
	fmt.Print(str)
	svc.outputFile.WriteString(str)
}
