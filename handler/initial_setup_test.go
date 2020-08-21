package handler

import "testing"

func TestSum(t *testing.T) {
	tables := []struct {
		machine string
		orders  string
	}{
		{"machine1", "orders1"},
	}
	for _, table := range tables {
		InitialSetup(table.machine, table.orders)
	}
}
