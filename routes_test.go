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

// register: /users/{:uid}
// /users/123 √
// /users/123/ √
// /users/123/info x
func Test_seekMatch(t *testing.T) {
	routesPath := "/users/{string:str}/{:defaultStr}/{int:num}/{float:floatNum}/{uint:uNum}"
	requestPath := "/users/clearString/defaultString/-12345/3.14159/12345"
	handler := func(writer http.ResponseWriter, request *http.Request) {
		params := Params(request)
		if params["str"].(string) != "clearString" {
			t.Fail()
		}
		if params["defaultStr"].(string) != "defaultString" {
			t.Fail()
		}
		if params["num"].(int) != -12345 {
			t.Fail()
		}
		if params["floatNum"].(float32) != 3.14159 {
			t.Fail()
		}
		if params["uNum"].(uint) != 12345 {
			t.Fail()
		}
		writer.WriteHeader(http.StatusOK)
	}

	r := NewRoutingTable()
	_, _ = r.Register(routesPath, handler, []string{http.MethodGet})
	wr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, requestPath, nil)
	r.ServeHTTP(wr, req)
	if wr.Code != http.StatusOK {
		t.Fail()
	}
}
