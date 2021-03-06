package go_restful_routes

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// RoutingTable
// When path does not contain a matching function, store it in fast(map);
// When path starts with `{`, store it in regex(slice);
// otherwise, store it in match(slice).
type RoutingTable struct {
	mux  *http.ServeMux
	full map[string]*routeItem // routesHash => *routeItem

	fast   map[string]*routeItem // `/users`
	prefix map[string]*routeItem // `/users/`
	match  []*routeItem          // `/users/{:id}/posts/{:post_id}`
	regex  []*routeItem          // `{^/[a-z]+\[[0-9]+\]$}`
	root   *routeItem            // `/`
}

func NewRoutingTable() *RoutingTable {
	return &RoutingTable{
		mux:    http.NewServeMux(),
		full:   make(map[string]*routeItem),
		fast:   make(map[string]*routeItem),
		prefix: make(map[string]*routeItem),
		regex:  make([]*routeItem, 0, 0),
		match:  make([]*routeItem, 0, 0),
	}
}

// ServeHTTP implemented `http.Handler` interface
func (r *RoutingTable) ServeHTTP(wrt http.ResponseWriter, req *http.Request) {
	originPath := req.URL.Path
	routeItem := r.seek(req) // Core
	Log(fmt.Sprintf("%v %v %v", req.Method, req.URL.Path, originPath))
	if routeItem != nil {
		handler, _ := r.mux.Handler(req)
		req.URL.Path = originPath
		ctx := context.WithValue(req.Context(), RouteItemContextKey, routeItem)
		handler.ServeHTTP(wrt, req.WithContext(ctx)) // handle by http.ServeMux
	} else {
		Log("route not found")
		http.NotFound(wrt, req)
	}
}

func (r *RoutingTable) Register(path string, handler func(http.ResponseWriter, *http.Request), methods []string) (item *routeItem, err error) {
	if item, err = newRouteItem(path, handler, methods); err != nil {
		return nil, err
	}
	switch {
	case path == "/":
		r.root = item
		r.mux.HandleFunc("/", handler)
	case path[0:1] == "{" && path[len(path)-1:] == "}":
		item.regex = path[1 : len(path)-1]
		r.regex = append(r.regex, item)
	case strings.Contains(path, "{"):
		if err = item.parsePatternBlocks(); err != nil {
			return nil, err
		}
		r.match = append(r.match, item)
	case path[len(path)-1:] == "/":
		r.prefix[path] = item
	default:
		r.fast[path] = item
	}
	r.full[item.key] = item
	r.mux.HandleFunc(item.key, handler)
	return item, nil
}

func (r *RoutingTable) Any(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, err := r.Register(path, handler, []string{}); err != nil {
		panic(err)
	}
}

func (r *RoutingTable) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, err := r.Register(path, handler, []string{http.MethodGet}); err != nil {
		panic(err)
	}
}

func (r *RoutingTable) Post(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, err := r.Register(path, handler, []string{http.MethodPost}); err != nil {
		panic(err)
	}
}

func (r *RoutingTable) Patch(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, err := r.Register(path, handler, []string{http.MethodPatch}); err != nil {
		panic(err)
	}
}

func (r *RoutingTable) Put(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, err := r.Register(path, handler, []string{http.MethodPut}); err != nil {
		panic(err)
	}
}

func (r *RoutingTable) Delete(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, err := r.Register(path, handler, []string{http.MethodDelete}); err != nil {
		panic(err)
	}
}

func (r *RoutingTable) Head(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, err := r.Register(path, handler, []string{http.MethodHead}); err != nil {
		panic(err)
	}
}

func (r *RoutingTable) Options(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, err := r.Register(path, handler, []string{http.MethodOptions}); err != nil {
		panic(err)
	}
}

func (r *RoutingTable) Connect(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, err := r.Register(path, handler, []string{http.MethodConnect}); err != nil {
		panic(err)
	}
}

func (r *RoutingTable) Trace(path string, handler func(http.ResponseWriter, *http.Request)) {
	if _, err := r.Register(path, handler, []string{http.MethodTrace}); err != nil {
		panic(err)
	}
}
