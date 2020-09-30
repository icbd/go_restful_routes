package go_restful_routes

import (
	"log"
)

var Verbose bool = true
var Log = func(s string) {
	if Verbose {
		log.Println(s)
	}
}
