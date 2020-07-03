package bot

import "net/url"

type ApiMethod interface {
	requestParam() url.Values
	method() string
}

type GetMe struct {
}

func (f GetMe) requestParam() url.Values { return nil }

func (f GetMe) method() string { return "getMe" }
