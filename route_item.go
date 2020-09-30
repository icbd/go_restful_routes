package go_restful_routes

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// routeItem
// key is a Hash generated by Path and HTTPMethods.
// The key acts not only as a key for RoutingTable.full, but also as a routes key for http.ServeMux.
type routeItem struct {
	HTTPMethods []string // request.Method
	Handler     func(http.ResponseWriter, *http.Request)
	Path        string // raw routes path
	key         string
	pathBlocks  []string // split Path by `/`
	regex       string   // trim `{}` from path
}

// newRouteItem Include sorting the methods and filling the key.
func newRouteItem(path string, handler func(http.ResponseWriter, *http.Request), methods []string) (*routeItem, error) {
	if path == "" {
		return nil, errors.New("[WARN] routing pattern cannot be empty")
	}
	item := &routeItem{Path: path, Handler: handler, HTTPMethods: methods}
	item.fillKey()
	return item, nil
}

// fillKey Calculate hash from path and http methods. Assign hash value to `item.key`.
func (item *routeItem) fillKey() *routeItem {
	bytes := []byte(item.Path)
	if len(item.HTTPMethods) > 0 {
		sort.Strings(item.HTTPMethods)
		for _, m := range item.HTTPMethods {
			bytes = append(bytes, m...)
		}
	}
	item.key = fmt.Sprintf("/%x", md5.Sum(bytes))
	return item
}

// validHTTPMethod empty slice means allow all
func (item *routeItem) validHTTPMethod(method string) bool {
	if len(item.HTTPMethods) == 0 {
		return true
	}
	for _, m := range item.HTTPMethods {
		if method == m {
			return true
		}
	}
	return false
}

func (item *routeItem) parsePathBlocks() (err error) {
	path := strings.TrimRight(item.Path, "/")
	blocks := strings.Split(path, "/")
	if len(blocks) < 1 {
		return errors.New(fmt.Sprintf("[WARN] routing pattern is invalid: %v", path))
	}
	item.pathBlocks = blocks
	return nil
}
