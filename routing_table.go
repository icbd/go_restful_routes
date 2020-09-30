package go_restful_routes

import (
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
	regex  []*routeItem          // `{^/[a-z]+\[[0-9]+\]$}`
	match  []*routeItem          // `/users/{:id}/posts/{:post_id}`
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
	routeItem := r.seek(req)
	Log(fmt.Sprintf("%v %v %v", req.Method, req.URL.Path, originPath))
	if routeItem != nil {
		handler, _ := r.mux.Handler(req)
		req.URL.Path = originPath
		handler.ServeHTTP(wrt, req) // handle by http.ServeMux
	} else {
		Log("route not found")
		http.NotFound(wrt, req)
	}
}

func (r *RoutingTable) Register(path string, handler func(http.ResponseWriter, *http.Request), methods []string) (item *routeItem, err error) {
	if item, err = newRouteItem(path, handler, methods); err != nil {
		return
	}
	switch {
	case path[0:1] == "{" && path[len(path)-1:] == "}":
		item.regex = path[1 : len(path)-1]
		r.regex = append(r.regex, item)
	case strings.Contains(path, "{"):
		item.pathBlocks = strings.Split(item.Path, "/")
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
