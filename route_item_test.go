package go_restful_routes

import (
	"net/http"
	"testing"
)

func Test_fillKey(t *testing.T) {
	// hash key should be filled
	item := &routeItem{
		Path:        "/users",
		HTTPMethods: []string{http.MethodGet},
	}
	item.fillKey()
	if len(item.key) != 32 {
		t.Fail()
	}

	// the hash of two different items should be different
	other := &routeItem{
		Path:        "/users/",
		HTTPMethods: []string{http.MethodGet},
	}
	other.fillKey()
	if item.key == other.key {
		t.Fail()
	}
}

func TestNewRouteItem(t *testing.T) {
	handler := func(writer http.ResponseWriter, request *http.Request) {}
	var err error
	var item *routeItem

	// path should not be empty
	item, err = newRouteItem("", handler, []string{http.MethodGet})
	if err == nil || item != nil {
		t.Fail()
	}

	// methods should not be empty
	item, err = newRouteItem("/", handler, []string{})
	if err == nil || item != nil {
		t.Fail()
	}

	// success
	item, err = newRouteItem("/", handler, []string{http.MethodGet})
	if err != nil || item == nil {
		t.Fail()
	} else if item.key == "" {
		t.Fail()
	}
}

func Test_validHTTPMethod(t *testing.T) {
	item, _ := newRouteItem("/", fakeHandler, []string{http.MethodGet, http.MethodPost})
	if !item.validHTTPMethod(http.MethodPost) {
		t.Fail()
	}
	if item.validHTTPMethod(http.MethodDelete) {
		t.Fail()
	}
}
