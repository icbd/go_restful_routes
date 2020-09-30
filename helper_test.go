package go_restful_routes

import "net/http"

func fakeHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
