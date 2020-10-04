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
// dir: Dictionary for storing static resources
// prefix: Routes prefix. URL: /prefix/hi.html => File: {{dir}}/hi.html
// suffix: Supported file suffix. If it is empty, all formats are supported.
type option struct {
	dir    string   `default:"public/"`
	prefix string   `default:"/" hint:"/static/"`
	suffix []string `default:"[ico, jpg, jpeg, png, gif, webp, html, js, css, md]"`
}

func New(wr http.ResponseWriter, req *http.Request, Dir string, Prefix string, Suffix ...string) *option {
	var opt *option
	default_box.New(&opt).Fill()
	http.FileServer(opt).ServeHTTP(wr, req)
	return opt
}

// Open implement: FileSystem interface
func (opt *option) Open(filePath string) (http.File, error) {
	if go_restful_routes.Verbose {
		go_restful_routes.Log(fmt.Sprintf("server for static file: %v", filePath))
	}

	if name, err := opt.parse(filePath); err != nil {
		return nil, err
	} else {
		return http.Dir(opt.dir).Open(name)
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
	if opt.prefix == "/" {
		return filePath, nil
	}

	if !strings.HasPrefix(filePath, opt.prefix) {
		return "", errors.New("prefix not matching")
	}

	parsed = strings.Replace(filePath, opt.prefix, "/", 1)

	return parsed, nil
}

func (opt *option) parseSuffix(filePath string) (suffix string, err error) {
	for _, s := range opt.suffix {
		if strings.HasSuffix(filePath, "."+s) {
			return s, nil
		}
	}

	return "", errors.New(fmt.Sprintf("suffix (%s) not supported", filePath))
}
