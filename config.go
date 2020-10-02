package go_restful_routes

import (
	"log"
)

const RouteItemContextKey string = "RouteItemContextKey"

var Verbose bool = true
var Log = func(s string) {
	if Verbose {
		log.Println(s)
	}
}
