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
	values := url.Values{}
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

type SendMessage struct {
	ChatID                int64       `json:"chat_id"`
	Text                  string      `json:"text"`
	ParseMode             string      `json:"parse_mode"`
	DisableWebPagePreview bool        `json:"disable_web_page_preview"`
	DisableNotification   bool        `json:"disable_notification"`
	ReplyToMessageID      int         `json:"reply_to_message_id"`
	ReplyMarkup           interface{} `json:"reply_markup"`
}

func (m SendMessage) requestParam() url.Values {
	values := url.Values{}
	values.Add("chat_id", strconv.FormatInt(m.ChatID, 10))
	values.Add("text", m.Text)

	return values
}

func (m SendMessage) method() string { return "sendMessage" }
