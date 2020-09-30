package go_restful_routes

import (
	"net/http"
	"testing"
)

func TestNewRoutingTable(t *testing.T) {
	// should initialize each map and slice
	routingTable := NewRoutingTable()
	if routingTable.mux == nil {
		t.Fail()
	}
	if routingTable.full == nil {
		t.Fail()
	}
	if routingTable.fast == nil {
		t.Fail()
	}
	if routingTable.regex == nil {
		t.Fail()
	}
	if routingTable.match == nil {
		t.Fail()
	}
}

func TestRegister(t *testing.T) {
	handler := func(writer http.ResponseWriter, request *http.Request) {}
	var item *RouteItem
	var err error
	var table *RoutingTable

	// add simple path to fast group
	table = NewRoutingTable()
	item, err = table.Register("/", handler, []string{http.MethodGet})
	if err != nil || item == nil {
		t.Fail()
	}
	if len(table.full) != 1 && len(table.fast) != 1 {
		t.Fail()
	}

	// add regex path to regex group
	table = NewRoutingTable()
	item, err = table.Register("{^/[a-z]+\\[[0-9]+\\]$}", handler, []string{http.MethodGet})
	if err != nil || item == nil {
		t.Fail()
	}
	if len(table.full) != 1 && len(table.regex) != 1 {
		t.Fail()
	}

	// add params path to match group
	table = NewRoutingTable()
	item, err = table.Register("/users/{:userId}/info/", handler, []string{http.MethodGet})
	if err != nil || item == nil {
		t.Fail()
	}
	if len(table.full) != 1 && len(table.match) != 1 {
		t.Fail()
	}

	// ignore it when path or methods invalid
	table = NewRoutingTable()
	item, err = table.Register("", handler, []string{http.MethodGet})
	if err == nil || item != nil {
		t.Fail()
	}
	if len(table.full) != 0 {
		t.Fail()
	}
}
