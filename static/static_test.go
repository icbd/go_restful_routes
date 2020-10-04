package static

import (
	"testing"
)

func opt() *option {
	return &option{
		dir:    "public/",
		prefix: "/static/",
		suffix: []string{"ico"},
	}
}

func Test_parsePrefix(t *testing.T) {
	opt := opt()
	var parsed string
	var err error

	parsed, err = opt.parsePrefix("/static/pic.ico")
	if parsed != "/pic.ico" {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}

	parsed, err = opt.parsePrefix("/sth/pic.ico")
	if err == nil {
		t.Fail()
	}
}

func Test_parseSuffix(t *testing.T) {
	opt := opt()
	var err error

	_, err = opt.parseSuffix("/static/pic.ico")
	if err != nil {
		t.Fail()
	}

	_, err = opt.parseSuffix("/static/pic.png")
	if err == nil {
		t.Fail()
	}

}
