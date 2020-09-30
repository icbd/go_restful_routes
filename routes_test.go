package go_restful_routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_seekFast(t *testing.T) {
	r := NewRoutingTable()
	_, _ = r.Register("/users", fakeHandler, []string{http.MethodGet})

	// success
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusOK {
		t.Fail()
	}

	// auto match `/users/` => `/users`
	wr = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/", nil)
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusOK {
		t.Fail()
	}
}
