package handler

import (
	"coffeemachine"
	"coffeemachine/service"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
)

func InitialSetup(machineConfig, orderFileName string) error {
	machineFile, err := ioutil.ReadFile(machineConfig)
	if err != nil {
		return err
	}
	var machine coffeemachine.Machine
	err = json.Unmarshal([]byte(machineFile), &machine)
	if err != nil {
		return errors.New("Unavle to unmarshal " + machineConfig + err.Error())
	}

	// starts the machine service
	svc := service.New(machine, machineConfig+"_"+orderFileName)

	var orders map[coffeemachine.BeverageName]int
	ordersFile, err := ioutil.ReadFile(orderFileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(ordersFile), &orders)
	if err != nil {
		return errors.New("Unavle to unmarshal " + orderFileName)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var orderNo int64
	orderNo = 1

	svc.Menu()
	for bvg, quant := range orders {
		for i := 0; i < quant; i++ {
			svc.TakeOrders(ctx, bvg, orderNo)
		}
		orderNo += 1
	}

	svc.AllDone(ctx)

	return nil
}
