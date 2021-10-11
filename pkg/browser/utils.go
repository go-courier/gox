package browser

import (
	"net/url"
	"strings"
)

func EncodeURIComponent(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
}
