package go_restful_routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// register: /user
// when /user  √
// when /user/ x
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
	if wr.Code != http.StatusNotFound {
		t.Fail()
	}
}

// register: /user/
// when /user/ √
// when /user/hi/hello √
// when /user/hi/hello/ √
func Test_seekPrefix(t *testing.T) {
	r := NewRoutingTable()
	_, _ = r.Register("/users/", fakeHandler, []string{http.MethodGet})

	// success
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusOK {
		t.Fail()
	}

	wr = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/hi", nil)
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusOK {
		t.Fail()
	}

	wr = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/hi/", nil)
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusOK {
		t.Fail()
	}

	wr = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/hi/hello", nil)
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusOK {
		t.Fail()
	}

	wr = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/hi/hello/", nil)
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusOK {
		t.Fail()
	}
}

func Test_seekRegex(t *testing.T) {
	r := NewRoutingTable()
	_, _ = r.Register("{^/user\\[[0-9]+\\]$}", fakeHandler, []string{http.MethodGet})

	// success
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/user[123]", nil)
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusOK {
		t.Fail()
	}

	// not match
	wr = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/user[]", nil)
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusNotFound {
		t.Fail()
	}
}
