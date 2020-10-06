package main

import (
	"encoding/json"
	"github.com/icbd/go_restful_routes"
	"github.com/icbd/go_restful_routes/static"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := http.ListenAndServe(":3000", Handler()); err != nil {
		panic(err)
	}
}

func Handler() http.Handler {
	r := go_restful_routes.NewRoutingTable()
	r.Get("/hi", HiController)
	r.Post("/users/{int:Uid}", ShowUser)
	r.Any("/", RootController)
	return r
}

type user struct {
	Uid       int       `json:"uid"`
	CreatedAt time.Time `json:"created_at"`
}

// GET /hi
func HiController(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write([]byte("<h1>hi</h1>"))
}

// GET /users/123
func ShowUser(w http.ResponseWriter, req *http.Request) {
	params := go_restful_routes.Params(req)
	u := user{Uid: params["Uid"].(int), CreatedAt: time.Now()}
	if data, err := json.Marshal(u); err == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

// ANY /
// ANY /404.html
// ANY /avatar.png
func RootController(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		static.New(w, req, os.Getenv("PUBLIC_DIR"), "/", "png", "html")
		return
	}
	_, _ = w.Write([]byte("<h1>Welcome~</h1>"))
}
