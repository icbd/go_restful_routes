package static

import (
	"errors"
	"fmt"
	"github.com/icbd/default_box"
	"github.com/icbd/go_restful_routes"
	"net/http"
	"os"
	"strings"
)

// option
// Dir: Dictionary for storing static resources
// Prefix: Routes Prefix. URL: /Prefix/hi.html => File: {{Dir}}/hi.html
// Suffix: Supported file Suffix. If it is empty, all formats are supported.
type option struct {
	Dir    string   `default:"public/"`
	Prefix string   `default:"/" hint:"/static/"`
	Suffix []string `default:"[ico, jpg, jpeg, png, gif, webp, html, js, css, md]"`
}

func New(wr http.ResponseWriter, req *http.Request, Dir string, Prefix string, Suffix ...string) *option {
	var opt option
	default_box.New(&opt).Fill()
	opt.Dir = Dir
	opt.Prefix = Prefix
	if len(Suffix) == 0 {
		opt.Suffix = Suffix
	}
	http.FileServer(&opt).ServeHTTP(wr, req)
	return &opt
}

// Open implement: FileSystem interface
func (opt *option) Open(filePath string) (http.File, error) {
	if go_restful_routes.Verbose {
		go_restful_routes.Log(fmt.Sprintf("server for static file: %v", filePath))
	}

	if name, err := opt.parse(filePath); err != nil {
		return nil, err
	} else {
		return http.Dir(opt.Dir).Open(name)
	}
}

func (opt *option) parse(filePath string) (parsed string, err error) {
	if _, err := opt.parseSuffix(filePath); err != nil {
		if go_restful_routes.Verbose {
			go_restful_routes.Log(err.Error())
		}
		return "", os.ErrNotExist
	}

	if parsed, err = opt.parsePrefix(filePath); err != nil {
		if go_restful_routes.Verbose {
			go_restful_routes.Log(fmt.Sprintln(err))
		}
		return "", os.ErrNotExist
	}

	return parsed, nil
}

func (opt *option) parsePrefix(filePath string) (parsed string, err error) {
	if opt.Prefix == "/" {
		return filePath, nil
	}

	if !strings.HasPrefix(filePath, opt.Prefix) {
		return "", errors.New("Prefix not matching")
	}

	parsed = strings.Replace(filePath, opt.Prefix, "/", 1)

	return parsed, nil
}

func (opt *option) parseSuffix(filePath string) (suffix string, err error) {
	for _, s := range opt.Suffix {
		if strings.HasSuffix(filePath, "."+s) {
			return s, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Suffix (%s) not supported", filePath))
}
