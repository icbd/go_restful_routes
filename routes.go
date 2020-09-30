package go_restful_routes

import (
	"net/http"
	"regexp"
	"strings"
)

func (r *RoutingTable) seek(req *http.Request) (item *routeItem) {
	item = r.seekFast(req)
	if item != nil {
		return item
	}

	item = r.seekPrefix(req)
	if item != nil {
		return item
	}

	item = r.seekRegex(req)
	if item != nil {
		return item
	}

	return nil
}

func (r *RoutingTable) seekFast(req *http.Request) (item *routeItem) {
	if item, ok := r.fast[req.URL.Path]; ok {
		if item.validHTTPMethod(req.Method) {
			req.URL.Path = item.key
			return item
		}
	}
	return nil
}

func (r *RoutingTable) seekPrefix(req *http.Request) (item *routeItem) {
	for _, item = range r.prefix {
		if strings.HasPrefix(req.URL.Path, item.Path) {
			if item.validHTTPMethod(req.Method) {
				req.URL.Path = item.key
				return item
			}
		}
	}
	return nil
}

func (r *RoutingTable) seekRegex(req *http.Request) (item *routeItem) {
	for _, item = range r.regex {
		if matched, err := regexp.MatchString(item.regex, req.URL.Path); matched && err == nil {
			if item.validHTTPMethod(req.Method) {
				req.URL.Path = item.key
				return item
			}
		}
	}
	return nil
}

//func (r *RoutingTable) seekMatch(req *http.Request) (item *routeItem) {
//
//}
