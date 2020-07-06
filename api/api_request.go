package api

import (
	"fmt"
	"net/url"
	"strconv"
)

type ApiMethod interface {
	requestParam() url.Values
	method() string
}

type GetMe struct {
}

func (m GetMe) requestParam() url.Values { return nil }

func (m GetMe) method() string { return "getMe" }

func (m GetMe) String() string { return "getMe\n" }

type GetUpdates struct {
	Offset  int
	Limit   int
	Timeout int
}

func (m GetUpdates) requestParam() url.Values {
	values := make(url.Values)
	if m.Offset != 0 {
		values.Add("offset", strconv.Itoa(m.Offset))
	}

	if m.Limit > 0 {
		values.Add("limit", strconv.Itoa(m.Limit))
	}

	if m.Timeout > 0 {
		values.Add("timeout", strconv.Itoa(m.Timeout))
	}

	return values
}

func (m GetUpdates) method() string { return "getUpdates" }

func (m GetUpdates) String() string {
	return fmt.Sprintf("Offset: %d, Limit: %d, Timeout: %d\n", m.Offset, m.Limit, m.Timeout)
}
