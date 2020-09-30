package go_restful_routes

import (
	"net/http"
)

func (r *RoutingTable) seek(req *http.Request) (item *routeItem) {
	item = r.seekFast(req)
	if item != nil {
		req.URL.Path = item.key

	}

	return
}

func (r *RoutingTable) seekFast(req *http.Request) *routeItem {
	path := normalizeFastPath(req.URL.Path)
	if item, ok := r.fast[path]; ok {
		if item.validHTTPMethod(req.Method) {
			return item
		}
	}
	return nil
}
