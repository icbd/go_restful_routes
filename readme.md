# Go RESTFul Routes

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/icbd/go_restful_routes/Test)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/icbd/go_restful_routes)

`go_restful_routes` is designed to build RESTFul routes in a concise way.

## Install

```shell script
go get github.com/icbd/go_restful_routes
```

## Run Example

```shell script
$ PUBLIC_DIR=example/public go run example/main.go
```

## How to use

Initialize a new routing table using `NewRoutingTable`. 
Registering routing entries by `Register`.

The first parameter of `Register` is the matching path, it has 5 types, and will be parsed in order:

1. fast matching, exact string matching. ex, `/users` can be matched by `/users`;
2. prefix matching. ex, `/users/`, `/users/123`, `/users/123/info`, `/users/123/info/` can matched by `/users/`;
3. params matching. ex, `/users/123` can matched by `/users/{int:uid}`, `/users/bob` can matched by `/users/{:name}`;
4. regex matching. ex, `/users[123]` can matched by `{^/users\[[0-9]+\]$}`;
5. root matching. Only match `/` .    

The twice parameter of `Register` is the handler function.

The third parameter of `Register` is the slice of http methods. Empty slice represent any method allowed.

```go
func Handler() http.Handler {
	r := go_restful_routes.NewRoutingTable()
	_, _ = r.Register("/", controllers.RootController, []string{})
	_, _ = r.Register("/hi", controllers.HiController, []string{http.MethodGet, http.MethodPost})
	_, _ = r.Register("/users/{int:uid}", controllers.ShowUser, []string{http.MethodGet})
	return r
}
```

You can also use HTTP method verb to register routes directly.

```go
r.Any("/", controllers.RootController)
r.Get("/hi", controllers.HiController)
```

## Static Routes

```go
func RootController(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		static.New(writer, request, "public/", "/static/", "ico", "html")
		return
	}

	_, _ = writer.Write([]byte("<h1>Welcome~</h1>"))
}

```

* Using `Dir` to specify static resource directory;

* Using `Prefix` to replace routes prefix, ex, if request URL is`/static/hello.png`, will Server the file in `{Dir}/hello.png`;

* Using `Suffix` to specify static file types. If leave it blank, the default value is used: `[ico, jpg, jpeg, png, gif, webp, html, js, css, md]`

## Config Log

```go
go_restful_routes.Verbose = true
go_restful_routes.Log = func(s string) {
    if go_restful_routes.Verbose {
        log.Println(s)
    }
}
```

## License

MIT, see [LICENSE](LICENSE)
